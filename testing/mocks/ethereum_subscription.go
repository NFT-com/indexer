package mocks

import (
	"testing"
)

type Subscription struct {
	UnsubscribeFunc func()
	ErrFunc         func() <-chan error
}

func BaselineSubscription(t *testing.T) *Subscription {
	t.Helper()

	c := Subscription{
		UnsubscribeFunc: func() {},
		ErrFunc: func() <-chan error {
			return GenericErrChannel
		},
	}

	return &c
}

func (c *Subscription) Unsubscribe() {
	c.UnsubscribeFunc()
}

func (c *Subscription) Err() <-chan error {
	return c.ErrFunc()
}
