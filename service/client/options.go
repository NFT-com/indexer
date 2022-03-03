package client

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type options struct {
	websocketURL url.URL
	httpURL      url.URL
	httpClient   *http.Client
	wsDialer     *websocket.Dialer
}

type Option func(*options)

type OptionsList []Option

func NewOptions(opts ...Option) OptionsList {
	return opts
}

func defaultOptions() *options {
	return &options{
		websocketURL: url.URL{Scheme: "ws", Host: "localhost:8081"},
		httpURL:      url.URL{Scheme: "http", Host: "localhost:8081"},
		httpClient:   http.DefaultClient,
		wsDialer:     websocket.DefaultDialer,
	}
}

func WithWebsocketScheme(scheme string) Option {
	return func(o *options) {
		o.websocketURL.Scheme = scheme
	}
}

func WithHost(host string) Option {
	return func(o *options) {
		o.websocketURL.Host = host
		o.httpURL.Host = host
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *options) {
		o.httpClient = httpClient
	}
}

func (o OptionsList) Apply(opts *options) {
	for _, opt := range o {
		opt(opts)
	}
}
