package queue

import (
	"context"
	"sync"
	"time"
)

type canceller struct {
	sync.Mutex

	subscribers map[chan struct{}]string
	cancelled   map[string]time.Time
}

func newCanceller() *canceller {
	return &canceller{
		subscribers: make(map[chan struct{}]string),
		cancelled:   make(map[string]time.Time),
	}
}

func (c *canceller) Cancel(ctx context.Context, id string) error {
	c.Lock()
	c.cancelled[id] = time.Now().Add(time.Minute * 5)
	for subscriber, build := range c.subscribers {
		if id == build {
			close(subscriber)
		}
	}
	c.collect()
	c.Unlock()
	return nil
}

func (c *canceller) Cancelled(ctx context.Context, id string) (bool, error) {
	subscriber := make(chan struct{})
	c.Lock()
	c.subscribers[subscriber] = id
	c.Unlock()
	defer func() {
		c.Lock()
		delete(c.subscribers, subscriber)
		c.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(time.Minute):
			c.Lock()
			_, ok := c.cancelled[id]
			c.Unlock()
			if ok {
				return true, nil
			}
		case <-subscriber:
			return true, nil
		}
	}
}

func (c *canceller) collect() {
	// the list of cancelled builds is stored with a ttl, and
	// is not removed until the ttl is reached. This provides
	// adequate window for clients with connectivity issues to
	// reconnect and receive notification of cancel events.
	now := time.Now()
	for build, timestamp := range c.cancelled {
		if now.After(timestamp) {
			delete(c.cancelled, build)
		}
	}
}
