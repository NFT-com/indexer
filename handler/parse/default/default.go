package web3

import (
	"context"

	"github.com/rs/zerolog"
	
	"github.com/NFT-com/indexer/queue"
)

type Default struct {
	log zerolog.Logger
}

func NewDefault(log zerolog.Logger) (*Default, error) {
	d := Default{
		log: log.With().Str("component", "parse_default").Logger(),
	}

	return &d, nil
}

func (w *Default) Handle(ctx context.Context, job queue.ParseJob) error {
	w.log.Info().Interface("job", job).Msg("new job")
	return nil
}
