package web2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer res.Body.Close()

	_, ok := m.cfg.RetryCodes[res.StatusCode]
	if ok {
		return nil, fmt.Errorf("bad response code (%d)", res.StatusCode)
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return nil, fmt.Errorf("invalid response code (%d)", res.StatusCode)
	}

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	var token metadata.Token
	err = json.Unmarshal(payload, &token)
	if err != nil {
		return nil, fmt.Errorf("could not decode token metadata: %w", err)
	}

	return &token, nil
}
