package retry

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func Indefinite() *backoff.ExponentialBackOff {

	exp := backoff.ExponentialBackOff{
		InitialInterval:     200 * time.Millisecond,
		RandomizationFactor: 0.5,
		Multiplier:          2,
		MaxInterval:         1 * time.Minute,
		MaxElapsedTime:      0,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}

	return &exp
}
