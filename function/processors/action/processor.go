package action

import (
	"context"

	"github.com/NFT-com/indexer/jobs"
	"github.com/NFT-com/indexer/models/chain"
)

type Processor interface {
	Type() string
	Standard() string
	Process(ctx context.Context, job jobs.Action) (*chain.NFT, error)
}
