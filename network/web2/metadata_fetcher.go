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
}

func NewMetadataFetcher() *MetadataFetcher {

	m := MetadataFetcher{}

	return &m
}

func (m *MetadataFetcher) Token(ctx context.Context, uri string) (*metadata.Token, error) {

	res, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer res.Body.Close()

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
