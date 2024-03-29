package web2

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

type MetadataFetcher struct {
	client *http.Client
	cfg    MetadataConfig
}

func NewMetadataFetcher(options ...MetadataOption) *MetadataFetcher {

	cfg := MetadataDefaultConfig
	for _, option := range options {
		option(&cfg)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.DisableValidation},
	}
	client := &http.Client{Transport: transport}

	m := MetadataFetcher{
		client: client,
		cfg:    cfg,
	}

	return &m
}

func (m *MetadataFetcher) Payload(_ context.Context, uri string) ([]byte, error) {

	res, err := m.client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("bad response code (%d)", res.StatusCode)
	}

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	return payload, nil
}
