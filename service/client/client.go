package client

import (
	"github.com/rs/zerolog"
)

const (
	discoveryBasePath = "discoveries"
	parsingBasePath   = "parsings"

	contentTypeHeaderName = "content-type"
	jsonContentType       = "application/json"
)

type Client struct {
	log     zerolog.Logger
	options *options
	close   chan struct{}
}

func New(log zerolog.Logger, optionList OptionsList) *Client {
	opts := defaultOptions()
	optionList.Apply(opts)

	c := Client{
		log:     log.With().Str("component", "api_client").Logger(),
		options: opts,
		close:   make(chan struct{}),
	}

	return &c
}

func (c *Client) Close() {
	close(c.close)
}
