package processors

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"

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
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s-%s-%s", job.ChainID, job.Address, job.TokenID)))
	nftID := common.Bytes2Hex(hash[:])

	nft := chain.NFT{
		ID:       nftID,
		ChainID:  job.ChainID,
		Contract: job.Address,
		TokenID:  job.TokenID,
		Owner:    job.ToAddress,
	}

	return &nft, nil
}
