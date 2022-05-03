package processors

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"

	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/inputs"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/metadata"
)

type Metadata struct {
	client ethclient.Client
}

func NewMetadata(client ethclient.Client) *Metadata {

	m := Metadata{
		client: client,
	}

	return &m
}

func (m *Metadata) ERC721(ctx context.Context, action *jobs.Action) (*graph.NFT, error) {
	return m.fetch(ctx, action, abis.ERC721, fetches.ERC721)
}

func (m *Metadata) ERC1155(ctx context.Context, action *jobs.Action) (*graph.NFT, error) {
	return m.fetch(ctx, action, abis.ERC1155, fetches.ERC1155)
}

func (m *Metadata) fetch(ctx context.Context, action *jobs.Action, abi abi.ABI, fetch string) (*graph.NFT, error) {

	tokenID, ok := big.NewInt(0).SetString(action.TokenID, 0)
	if !ok {
		return nil, fmt.Errorf("could not convert token ID to big integer")
	}

	input, err := abi.Pack(fetch, tokenID)
	if err != nil {
		return nil, fmt.Errorf("could not pack input: %w", err)
	}

	address := common.HexToAddress(action.Address)
	msg := ethereum.CallMsg{From: common.Address{}, To: &address, Data: input}
	output, err := m.client.CallContract(ctx, msg, big.NewInt(0).SetUint64(action.Height))
	if err != nil {
		return nil, fmt.Errorf("could not call contract: %w", err)
	}

	fields, err := abi.Unpack(fetch, output)
	if err != nil {
		return nil, fmt.Errorf("could not unpack output: %w", err)
	}

	if len(fields) != 1 {
		return nil, fmt.Errorf("could not get uri: output len is not 1")
	}

	uri, ok := fields[0].(string)
	if !ok {
		return nil, fmt.Errorf("could not get uri: output is not a string")
	}

	resolved, err := resolveURI(uri)
	if err != nil {
		return nil, fmt.Errorf("could not resolve URI: %w", err)
	}

	res, err := http.Get(resolved)
	if err != nil {
		return nil, fmt.Errorf("could not perform get request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("could not perform get request: unexpected status code %d", res.StatusCode)
	}

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var token metadata.Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal json response: %w", err)
	}

	nftHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%s-%s", action.ChainID, action.Address, action.TokenID)))
	nftID := hex.EncodeToString(nftHash[:])

	traits := make([]*graph.Trait, 0, len(token.Attributes))
	for i, att := range token.Attributes {
		traitHash := sha3.Sum256([]byte(fmt.Sprintf("%s-%s-%s-%s", action.ChainID, action.Address, action.TokenID, i)))
		traitID := hex.EncodeToString(traitHash[:])
		trait := graph.Trait{
			ID:    traitID,
			Name:  att.TraitType,
			Value: fmt.Sprint(att.Value),
			NftID: nftID,
		}
		traits = append(traits, &trait)
	}

	var inputs inputs.Addition
	err = json.Unmarshal(action.Data, &inputs)
	if err != nil {
		return nil, fmt.Errorf("could not decode inputs: %w", err)
	}

	nft := graph.NFT{
		ID:          nftID,
		ChainID:     action.ChainID,
		Contract:    action.Address,
		TokenID:     action.TokenID,
		Name:        token.Name,
		URI:         uri,
		Image:       token.Image,
		Description: token.Description,
		Owner:       inputs.ToAddress,
		Traits:      traits,
	}

	return &nft, nil
}
