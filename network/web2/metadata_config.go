package web2

import (
	"net/http"
	"time"
)

var MetadataDefaultConfig = MetadataConfig{
	RetryCodes: map[int]struct{}{
		http.StatusTooManyRequests:     {},
		http.StatusInternalServerError: {},
		http.StatusBadGateway:          {},
		http.StatusServiceUnavailable:  {},
		http.StatusGatewayTimeout:      {},
	},
	RetryCap: 15 * time.Minute,
}

type MetadataConfig struct {
	RetryCodes map[int]struct{}
	RetryCap   time.Duration
}

type MetadataOption func(*MetadataConfig)

func WithRetryCodes(codes ...int) MetadataOption {
	return func(cfg *MetadataConfig) {
		cfg.RetryCodes = make(map[int]struct{})
		for _, code := range codes {
			cfg.RetryCodes[code] = struct{}{}
		}
	}
}

func WithRetryCap(cap time.Duration) MetadataOption {
	return func(cfg *MetadataConfig) {
		cfg.RetryCap = cap
	}
}
