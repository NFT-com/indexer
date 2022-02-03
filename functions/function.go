package functions

import (
	"context"
	"strings"

	"github.com/NFT-com/indexer/event"
)

func Name(network, chain, custom string) string {
	if custom == "" {
		return network + chain
	}

	return network + chain + strings.ToLower(custom)
}

type Function func(ctx context.Context, event *event.Event) error
