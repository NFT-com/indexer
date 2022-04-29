package client

import (
	"github.com/rs/zerolog"
)

const (
	discoveryBasePath = "discoveries"
	parsingBasePath   = "parsings"
	actionBasePath    = "actions"

	contentTypeHeaderName = "content-type"
	jsonContentType       = "application/json"
)

type Client struct {
	log    zerolog.Logger
	config config
	close  chan struct{}
}

func New(log zerolog.Logger, opts ...Option) *Client {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	c := Client{
		log:    log.With().Str("component", "api_client").Logger(),
		config: cfg,
		close:  make(chan struct{}),
	}

	return &c
}

func (c *Client) Close() {
	close(c.close)
}
