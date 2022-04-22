package client

import (
	"net/http"
	"net/url"
)

var (
	defaultConfig = config{
		jobsAPI: url.URL{Scheme: "http", Host: "localhost:8081"},
		client:  http.DefaultClient,
	}
)

type config struct {
	jobsAPI url.URL
	client  *http.Client
}

type Option func(*config)

func WithHost(host string) Option {
	return func(o *config) {
		o.jobsAPI.Host = host
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(o *config) {
		o.client = httpClient
	}
}
