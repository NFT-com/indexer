package client

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	defaultConfig = config{
		jobsWebsocket: url.URL{Scheme: "ws", Host: "localhost:8081"},
		jobsAPI:       url.URL{Scheme: "http", Host: "localhost:8081"},
		client:        http.DefaultClient,
		dialer:        websocket.DefaultDialer,
	}
)

type config struct {
	jobsWebsocket url.URL
	jobsAPI       url.URL
	client        *http.Client
	dialer        *websocket.Dialer
}

type Option func(*config)

func WithWebsocketScheme(scheme string) Option {
	return func(o *config) {
		o.jobsWebsocket.Scheme = scheme
	}
}

func WithHost(host string) Option {
	return func(o *config) {
		o.jobsWebsocket.Host = host
		o.jobsAPI.Host = host
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *config) {
		o.client = httpClient
	}
}
