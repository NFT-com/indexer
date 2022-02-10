package subscriber_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/source"
	"github.com/NFT-com/indexer/subscriber"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewSubscriber(t *testing.T) {
	tts := []struct {
		name          string
		log           zerolog.Logger
		parser        block.Parser
		sources       []source.Source
		expectedError bool
	}{
		{
			name:   "should return error on missing parser",
			log:    zerolog.Logger{},
			parser: nil,
			sources: []source.Source{
				mocks.BaselineSource(t),
			},
			expectedError: true,
		},
		{
			name:          "should return error on missing sources",
			log:           zerolog.Logger{},
			parser:        mocks.BaselineParser(t),
			sources:       []source.Source{},
			expectedError: true,
		},
		{
			name:   "should return no error",
			log:    zerolog.Logger{},
			parser: mocks.BaselineParser(t),
			sources: []source.Source{
				mocks.BaselineSource(t),
			},
			expectedError: false,
		},
	}
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			subs, err := subscriber.NewSubscriber(tt.log, tt.parser, tt.sources)
			if tt.expectedError && err == nil {
				t.Errorf("test %s failed expected error but got none", tt.name)
				return
			}

			if !tt.expectedError && subs == nil {
				t.Errorf("test %s failed expected subscriber but found none", tt.name)
				return
			}
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	t.Run("should run with any problem and return", func(t *testing.T) {
		t.Parallel()

		var (
			ctx     = context.Background()
			events  = make(chan *event.Event)
			log     = zerolog.New(os.Stderr)
			parser  = mocks.BaselineParser(t)
			source1 = mocks.BaselineSource(t)
			source2 = mocks.BaselineSource(t)
			sources = []source.Source{source1, source2}
		)

		subs, err := subscriber.NewSubscriber(log, parser, sources)
		if err != nil {
			t.Errorf("unexpected error creating subscriber")
			return
		}

		var (
			index      = 0
			eventCount = 0
			block1     = block.Block("block_hash_1")
			block2     = block.Block("block_hash_2")
			block3     = block.Block("block_hash_3")
			blocks     = []*block.Block{
				&block1, &block2, nil,
				&block3, nil,
			}
		)

		source1.NextFunc = func(_ context.Context) *block.Block {
			b := blocks[index]
			index++
			return b
		}

		source2.NextFunc = func(_ context.Context) *block.Block {
			b := blocks[index]
			index++
			return b
		}

		parser.ParseFunc = func(_ context.Context, _ *block.Block) ([]*event.Event, error) {
			if index%2 == 0 {
				return nil, errors.New("failed to parse event")
			}

			return []*event.Event{{ID: "id"}}, nil
		}

		go func() {
			for {
				<-events
				eventCount++
			}
		}()

		err = subs.Subscribe(ctx, events)
		if err != nil {
			t.Errorf("unexpected error closing subscriber")
			return
		}

		if eventCount != 1 {
			t.Errorf("unexpected event count: expected %v got %v", 1, eventCount)
			return
		}
	})

	t.Run("should stop the subscription", func(t *testing.T) {
		t.Parallel()

		var (
			ctx     = context.Background()
			events  = make(chan *event.Event)
			log     = zerolog.New(os.Stderr)
			parser  = mocks.BaselineParser(t)
			s       = mocks.BaselineSource(t)
			sources = []source.Source{s}
		)

		subs, err := subscriber.NewSubscriber(log, parser, sources)
		if err != nil {
			t.Errorf("unexpected error creating subscriber")
			return
		}

		b := block.Block("block_hash_1")
		s.NextFunc = func(_ context.Context) *block.Block {
			return &b
		}

		parser.ParseFunc = func(_ context.Context, _ *block.Block) ([]*event.Event, error) {
			return []*event.Event{{ID: "id"}}, nil
		}

		go func() {
			for {
				<-events
			}
		}()

		done := make(chan struct{})
		go func(t *testing.T) {
			t.Helper()
			err = subs.Subscribe(ctx, events)
			if err != nil {
				t.Errorf("unexpected error closing subscriber")
				return
			}

			close(done)
		}(t)

		_ = subs.Close()
		select {
		case <-time.After(time.Second * 15):
			t.Fatal("timeout running test")
		case <-done:
			return
		}
	})
}

func TestSubscriber_Close(t *testing.T) {
	t.Run("should close all the sources on close", func(t *testing.T) {
		t.Parallel()

		var (
			log     = zerolog.New(os.Stderr)
			parser  = mocks.BaselineParser(t)
			source1 = mocks.BaselineSource(t)
			source2 = mocks.BaselineSource(t)
			sources = []source.Source{source1, source2}
		)

		subs, err := subscriber.NewSubscriber(log, parser, sources)
		if err != nil {
			t.Errorf("unexpected error creating subscriber")
			return
		}

		closed := 0
		source1.CloseFunc = func() error {
			closed++
			return nil
		}

		source2.CloseFunc = func() error {
			closed++
			return errors.New("failed to close source")
		}

		err = subs.Close()
		if err != nil {
			t.Errorf("unexpected error closing subscriber")
			return
		}

		if closed != len(sources) {
			t.Errorf("expected to close %v sources but only closed %v", len(sources), closed)
			return
		}
	})
}
