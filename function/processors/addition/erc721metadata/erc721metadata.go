package erc721metadata

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/networks"
)

type Processor struct {
	log     zerolog.Logger
	network networks.Network
	abi     abi.ABI
	client  *http.Client
}

func NewProcessor(log zerolog.Logger, network networks.Network) (*Processor, error) {
	abi, err := abi.JSON(strings.NewReader(uriFunctionABI))
	if err != nil {
		return nil, fmt.Errorf("could not parse abi: %w", err)
	}

	h := Processor{
		log:     log,
		network: network,
		abi:     abi,
		client:  http.DefaultClient,
	}

	return &h, nil
}

func (p *Processor) Process(ctx context.Context, job jobs.Addition) (*chain.NFT, error) {
	chainID, err := p.network.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get chain id: %w", err)
	}

	uri, err := p.getURI(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("could not get uri: %w", err)
	}

	resp, err := p.client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("could not perform get request: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("could not perform get request: unexpected status code %d", resp.StatusCode)
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var info TokenMetadata
	err = json.Unmarshal(payload, &info)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s", chainID, job.Address, job.TokenID)))
	nftID := common.Bytes2Hex(hash[:])

	traits := make([]chain.Trait, 0, len(info.Attributes))
	traitMap := make(map[string]int, len(info.Attributes))
	for _, att := range info.Attributes {
		traitHash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s-%s", chainID, job.Address, job.TokenID, att.TraitType)))
		traitID := common.Bytes2Hex(traitHash[:])

		count, ok := traitMap[traitID]
		if ok {
			traitID = fmt.Sprintf("%s-%d", traitID, count)
		}
		traitMap[traitID] = count + 1

		var traitValue string
		switch att.DisplayType {
		case dateDisplayType:
			val, ok := att.Value.(float64)
			if !ok {
				return nil, fmt.Errorf("could not parse date time: type is %T and not float64", att.Value)
			}

			traitValue = time.Unix(int64(val), 0).Format(time.RFC3339)
		default:
			traitValue = fmt.Sprint(att.Value)
		}

		traits = append(traits, chain.Trait{
			ID:    traitID,
			Name:  att.TraitType,
			Value: traitValue,
			NftID: nftID,
		})
	}

	nft := chain.NFT{
		ID:          nftID,
		ChainID:     chainID,
		Contract:    job.Address,
		TokenID:     job.TokenID,
		Name:        info.Name,
		Image:       info.Image,
		Description: info.Description,
		Traits:      traits,
	}

	return &nft, nil
}

func (p *Processor) getURI(ctx context.Context, job jobs.Addition) (string, error) {
	tokenID, ok := big.NewInt(0).SetString(job.TokenID, 0)
	if !ok {
		return "", fmt.Errorf("could not parse token id to big int")
	}

	input, err := p.abi.Pack(tokenURIFunctionName, tokenID)
	if err != nil {
		return "", fmt.Errorf("could not pack input: %w", err)
	}

	output, err := p.network.CallContract(ctx, nil, callSender, job.Address, input)
	if err != nil {
		return "", fmt.Errorf("could not call function on contract: %w", err)
	}

	outputFields, err := p.abi.Unpack(tokenURIFunctionName, output)
	if err != nil {
		return "", fmt.Errorf("could not unpack output: %w", err)
	}

	if len(outputFields) != 1 {
		return "", fmt.Errorf("could not get uri: output len is not 1")
	}

	uri, ok := outputFields[0].(string)
	if !ok {
		return "", fmt.Errorf("could not get uri: output is not a string")
	}

	return uri, nil
}
