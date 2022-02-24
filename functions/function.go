package functions

import (
	"context"
	"github.com/NFT-com/indexer/event"
)

func Name(network, chain string) string {
	return network + chain
}

type Function func(ctx context.Context, event *event.Event) error
