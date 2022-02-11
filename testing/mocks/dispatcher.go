package mocks

import (
	"context"
	"github.com/NFT-com/indexer/event"
	"testing"
)

type Dispatcher struct {
	DispatchFunc func(ctx context.Context, event *event.Event) error
}

func BaselineDispatcher(t *testing.T) *Dispatcher {
	t.Helper()

	c := Dispatcher{
		DispatchFunc: func(context.Context, *event.Event) error {
			return nil
		},
	}

	return &c
}

func (s *Dispatcher) Dispatch(ctx context.Context, event *event.Event) error {
	return s.DispatchFunc(ctx, event)
}
