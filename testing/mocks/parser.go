package mocks

import (
	"context"
	"testing"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
)

type Parser struct {
	ParseFunc func(ctx context.Context, block *block.Block) ([]*event.Event, error)
}

func BaselineParser(t *testing.T) *Parser {
	t.Helper()

	c := Parser{
		ParseFunc: func(ctx context.Context, block *block.Block) ([]*event.Event, error) {
			return GenericEvents, nil
		},
	}

	return &c
}

func (s *Parser) Parse(ctx context.Context, block *block.Block) ([]*event.Event, error) {
	return s.ParseFunc(ctx, block)
}
