package client

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

type options struct {
	jobsWebsocket url.URL
	jobsAPI       url.URL
	client        *http.Client
	dialer        *websocket.Dialer
}

type Option func(*options)

type OptionsList []Option

func NewOptions(opts ...Option) OptionsList {
	return opts
}

func defaultOptions() *options {
	return &options{
		jobsWebsocket: url.URL{Scheme: "ws", Host: "localhost:8081"},
		jobsAPI:       url.URL{Scheme: "http", Host: "localhost:8081"},
		client:        http.DefaultClient,
		dialer:        websocket.DefaultDialer,
	}
}

func WithWebsocketScheme(scheme string) Option {
	return func(o *options) {
		o.jobsWebsocket.Scheme = scheme
	}
}

func WithHost(host string) Option {
	return func(o *options) {
		o.jobsWebsocket.Host = host
		o.jobsAPI.Host = host
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *options) {
		o.client = httpClient
	}
}

func (o OptionsList) Apply(opts *options) {
	for _, opt := range o {
		opt(opts)
	}
}
