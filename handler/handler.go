package handler

import (
	"context"

	"github.com/NFT-com/indexer/queue"
)

type DiscoveryHandler interface {
	Handle(ctx context.Context, job queue.DiscoveryJob) error
}

type ParseHandler interface {
	Handle(ctx context.Context, job queue.ParseJob) error
}
