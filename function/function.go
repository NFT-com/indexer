package function

import (
	"context"

	"github.com/NFT-com/indexer/event"
)

func Name(network, chain, custom string) string {
	if custom == "" {
		return network + "-" + chain
	}

	return network + "-" + chain + "-" + custom
}

type Function func(ctx context.Context, event *event.Event) error
