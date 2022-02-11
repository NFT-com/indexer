package mocks

import (
	"context"
	"testing"

	"github.com/NFT-com/indexer/block"
)

type Source struct {
	NextFunc  func(ctx context.Context) *block.Block
	CloseFunc func() error
}

func BaselineSource(t *testing.T) *Source {
	t.Helper()

	c := Source{
		NextFunc: func(context.Context) *block.Block {
			return &GenericBlock
		},
		CloseFunc: func() error {
			return nil
		},
	}

	return &c
}

func (s *Source) Next(ctx context.Context) *block.Block {
	return s.NextFunc(ctx)
}

func (s *Source) Close() error {
	return s.CloseFunc()
}
