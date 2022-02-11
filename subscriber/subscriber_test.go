package subscriber_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
	"github.com/NFT-com/indexer/source"
	"github.com/NFT-com/indexer/subscriber"
	"github.com/NFT-com/indexer/testing/mocks"
)

func TestNewSubscriber(t *testing.T) {
	tests := []struct {
		name        string
		log         zerolog.Logger
		parser      block.Parser
		sources     []source.Source
		assertValue assert.ValueAssertionFunc
		assertError assert.ErrorAssertionFunc
	}{
		{
			name:   "should return error on missing parser",
			log:    zerolog.Logger{},
			parser: nil,
			sources: []source.Source{
				mocks.BaselineSource(t),
			},
			assertValue: assert.Nil,
			assertError: assert.Error,
		},
		{
			name:        "should return error on missing sources",
			log:         zerolog.Logger{},
			parser:      mocks.BaselineParser(t),
			sources:     []source.Source{},
			assertValue: assert.Nil,
			assertError: assert.Error,
		},
		{
			name:   "should return no error",
			log:    zerolog.Logger{},
			parser: mocks.BaselineParser(t),
			sources: []source.Source{
				mocks.BaselineSource(t),
			},
			assertValue: assert.NotNil,
			assertError: assert.NoError,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			
			subs, err := subscriber.NewSubscriber(test.log, test.parser, test.sources)
			test.assertError(t, err)
			test.assertValue(t, subs)
		})
	}
}

func TestSubscriber_Subscribe(t *testing.T) {
	var (
		ctx        = context.Background()
		events     = make(chan *event.Event)
		log        = zerolog.New(os.Stderr)
		parser     = mocks.BaselineParser(t)
		source1    = mocks.BaselineSource(t)
		source2    = mocks.BaselineSource(t)
		sources    = []source.Source{source1, source2}
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

	t.Run("should run with any problem and return", func(t *testing.T) {
		subs, err := subscriber.NewSubscriber(log, parser, sources)
		require.NoError(t, err)

		index = 0
		source1Count := 0
		source2Count := 0
		source1.NextFunc = func(context.Context) *block.Block {
			b := blocks[index]
			index++
			source1Count++
			return b
		}

		source2.NextFunc = func(context.Context) *block.Block {
			b := blocks[index]
			index++
			source2Count++
			return b
		}

		parser.ParseFunc = func(context.Context, *block.Block) ([]*event.Event, error) {
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
		assert.NoError(t, err)
		assert.Equal(t, 3, source1Count)
		assert.Equal(t, 2, source2Count)
		assert.Equal(t, 1, eventCount)
	})

	t.Run("should stop the subscription", func(t *testing.T) {
		s := mocks.BaselineSource(t)
		subs, err := subscriber.NewSubscriber(log, parser, []source.Source{s})
		require.NoError(t, err)

		s.NextFunc = func(_ context.Context) *block.Block {
			return &block1
		}

		parser.ParseFunc = func(context.Context, *block.Block) ([]*event.Event, error) {
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

			require.NoError(t, subs.Subscribe(ctx, events))
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
	var (
		log     = zerolog.New(os.Stderr)
		parser  = mocks.BaselineParser(t)
		source1 = mocks.BaselineSource(t)
		source2 = mocks.BaselineSource(t)
		sources = []source.Source{source1, source2}
	)

	t.Run("should close all the sources on close", func(t *testing.T) {
		subs, err := subscriber.NewSubscriber(log, parser, sources)
		require.NoError(t, err)

		var (
			source1Closed bool
			source2Closed bool
		)
		source1.CloseFunc = func() error {
			source1Closed = true
			return nil
		}

		source2.CloseFunc = func() error {
			source2Closed = true
			return mocks.GenericError
		}

		assert.NoError(t, subs.Close())
		assert.True(t, source1Closed)
		assert.True(t, source2Closed)
	})
}
