package notifier

import (
	"time"
)

var DefaultConfig = Config{
	NotifyInterval: time.Second,
}

type Config struct {
	NotifyInterval time.Duration
}

type Option func(*Config)

func WithNotifyInterval(interval time.Duration) Option {
	return func(cfg *Config) {
		cfg.NotifyInterval = interval
	}
}
