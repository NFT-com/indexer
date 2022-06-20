package retry

import (
	"math"
	"time"

	"github.com/nsqio/go-nsq"
)

func Pipeline(inflight uint) *nsq.Config {

	cfg := nsq.NewConfig()

	cfg.ReadTimeout = 5 * time.Second
	// cfg.WriteTimeout = 1 * time.Second

	cfg.LookupdPollInterval = time.Second
	// cfg.LookupdPollJitter = 0.3
	// cfg.LookupdPollTimeout = 10 * time.Second

	cfg.MaxRequeueDelay = 5 * time.Minute
	cfg.DefaultRequeueDelay = time.Second

	// cfg.BackoffStrategy = &nsq.ExponentialStrategy{}
	cfg.MaxBackoffDuration = 2 * time.Minute
	// cfg.BackoffMultiplier = time.Second

	cfg.MaxAttempts = math.MaxUint16

	cfg.HeartbeatInterval = time.Second
	// cfg.SampleRate = 0

	cfg.Deflate = false
	// cfg.DeflateLevel = 0
	cfg.Snappy = true

	cfg.MaxInFlight = int(inflight)

	cfg.MsgTimeout = 15 * time.Minute

	return cfg
}
