package owner_change

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
	"github.com/NFT-com/indexer/networks"
)

const (
	hexadecimalBase = 16
)

var (
	errNoOwnerFound = errors.New("no owner found")
)

type Processor struct {
	log     zerolog.Logger
	network networks.Network
}

func NewProcessor(log zerolog.Logger, network networks.Network) (*Processor, error) {
	h := Processor{
		log:     log,
		network: network,
	}

	return &h, nil
}

func (p *Processor) Type() string {
	return processorType
}

func (p *Processor) Standard() string {
	return processorStandard
}

func (p *Processor) Process(ctx context.Context, job jobs.Action) (*chain.NFT, error) {
	chainID, err := p.network.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get chain id: %w", err)
	}

	tokenID, ok := big.NewInt(0).SetString(job.TokenID, 0)
	if !ok {
		return nil, fmt.Errorf("could not parse token id to big int")
	}

	owner, err := p.getOwner(ctx, tokenID, job)
	if err != nil {
		return nil, fmt.Errorf("could not get owner: %w", err)
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s", chainID, job.Address, job.TokenID)))
	nftID := common.Bytes2Hex(hash[:])

	nft := chain.NFT{
		ID:       nftID,
		ChainID:  chainID,
		Contract: job.Address,
		TokenID:  job.TokenID,
		Owner:    owner,
	}

	return &nft, nil
}

func (p *Processor) getOwner(ctx context.Context, nftID *big.Int, job jobs.Action) (string, error) {
	logs, err := p.network.BlockEvents(ctx, job.BlockNumber, job.BlockNumber, []string{job.Event}, []string{job.Address})
	if err != nil {
		return "", fmt.Errorf("could not get block events: %w", err)
	}

	var owner string
	for _, log := range logs {
		if len(log.IndexData) != defaultIndexDataLen {
			return "", fmt.Errorf("unexpected index data length (have: %d, want: %d)", len(log.IndexData), defaultIndexDataLen)
		}

		// If it is not the same NFT, ignore it.
		if common.HexToHash(nftID.Text(hexadecimalBase)) != common.HexToHash(log.IndexData[2]) {
			continue
		}

		owner = common.HexToAddress(log.IndexData[1]).String()
	}

	if owner == "" {
		return "", errNoOwnerFound
	}

	return owner, nil
}
