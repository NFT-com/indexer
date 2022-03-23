package client

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

const (
	discoveryBasePath = "discoveries"
	parsingBasePath   = "parsings"

	contentTypeHeaderName = "content-type"
	jsonContentType       = "application/json"
)

type Client struct {
	log        zerolog.Logger
	wsClient   *websocket.Dialer
	httpClient *http.Client
	options    *options
	close      chan struct{}
}

func NewClient(log zerolog.Logger, optionList OptionsList) *Client {
	opts := defaultOptions()
	optionList.Apply(opts)

	c := Client{
		log:        log.With().Str("component", "api_client").Logger(),
		wsClient:   opts.wsDialer,
		httpClient: opts.httpClient,
		options:    opts,
		close:      make(chan struct{}),
	}

	return &c
}

func (c *Client) Close() {
	close(c.close)
}
