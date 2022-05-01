package client

import (
	"fmt"
	"net/http"
	"net/url"
)

const (
	parsingBasePath = "parsings"
	actionBasePath  = "actions"

	contentTypeHeaderName = "content-type"
	jsonContentType       = "application/json"
)

type Client struct {
	url    *url.URL
	client *http.Client
}

func New(opts ...Option) (*Client, error) {

	cfg := DefaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	url, err := url.Parse(cfg.APIURL)
	if err != nil {
		return nil, fmt.Errorf("could not parse API URL: %w", err)
	}

	c := Client{
		url:    url,
		client: http.DefaultClient,
	}

	return &c, nil
}
