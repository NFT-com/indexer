package web2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"

	"github.com/NFT-com/indexer/config/retry"
	"github.com/NFT-com/indexer/models/metadata"
)

type MetadataFetcher struct {
	cfg MetadataConfig
}

func NewMetadataFetcher(options ...MetadataOption) *MetadataFetcher {

	cfg := MetadataDefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	m := MetadataFetcher{
		cfg: cfg,
	}

	return &m
}

func (m *MetadataFetcher) Token(_ context.Context, uri string) (*metadata.Token, error) {

	notify := func(err error, dur time.Duration) {
		log.Warn().Err(err).Dur("duration", dur).Msg("could not get token data, retrying")
	}

	var payload []byte
	err := backoff.RetryNotify(func() error {

		res, err := http.Get(uri)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not execute request: %w", err))
		}
		defer res.Body.Close()

		_, ok := m.cfg.RetryCodes[res.StatusCode]
		if ok {
			return fmt.Errorf("bad response code (%d)", res.StatusCode)
		}

		if res.StatusCode < 200 || res.StatusCode > 299 {
			return backoff.Permanent(fmt.Errorf("fatal response code (%d)", res.StatusCode))
		}

		payload, err = io.ReadAll(res.Body)
		if err != nil {
			return backoff.Permanent(fmt.Errorf("could not read response body: %w", err))
		}

		return nil
	}, retry.Capped(m.cfg.RetryCap), notify)
	if err != nil {
		return nil, fmt.Errorf("could not get token data: %w", err)
	}

	var token metadata.Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return nil, fmt.Errorf("could not decode token metadata: %w", err)
	}

	return &token, nil
}
