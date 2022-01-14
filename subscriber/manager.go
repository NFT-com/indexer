package subscriber

import (
	"context"

	"github.com/NFT-com/indexer/events"
)

type Manager struct {
	subscribers []Subscriber
}

func (m *Manager) Subscribe(ctx context.Context, events chan events.Event) error {
	// Pass in to subscribers
	return nil
}

func (m *Manager) Status(ctx context.Context) error {
	// Check status of every subscriber
	return nil
}

func (m *Manager) Close() error {
	// Close every subscriber
	return nil
}
