package queue

import (
	"context"
	"dubhe-ci/core"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func newQueue(store core.BuildStore) *queue {
	q := &queue{
		interval: 1,
		ready:    make(chan struct{}, 1),
		paused:   false,
		store:    store,
		ctx:      context.Background(),
		workers:  map[*worker]struct{}{},
	}

	go q.start()

	return q
}

type queue struct {
	sync.Mutex

	interval time.Duration
	ready    chan struct{}
	paused   bool
	store    core.BuildStore
	ctx      context.Context
	workers  map[*worker]struct{}
}

func (q *queue) start() error {
	logrus.WithField("scheduler", "queue").Infoln("starting...")
	for {
		select {
		case <-q.ctx.Done():
			return q.ctx.Err()
		case <-q.ready:
			_ = q.signal(q.ctx)
		case <-time.After(q.interval * time.Second):
			_ = q.signal(q.ctx)
		}
	}
}

type worker struct {
	channel chan *core.Build
}

func (q *queue) Request(ctx context.Context) (*core.Build, error) {
	w := &worker{channel: make(chan *core.Build)}
	q.Lock()
	q.workers[w] = struct{}{}
	q.Unlock()

	select {
	case q.ready <- struct{}{}:
	default:
	}

	select {
	case <-ctx.Done():
		q.Lock()
		delete(q.workers, w)
		q.Unlock()
		return nil, ctx.Err()
	case b := <-w.channel:
		return b, nil
	}
}

func (q *queue) Schedule(ctx context.Context, build *core.Build) error {
	select {
	case q.ready <- struct{}{}:
	default:
	}
	return nil
}

func (q *queue) Pause(ctx context.Context) error {
	q.Lock()
	q.paused = true
	q.Unlock()
	return nil
}

func (q *queue) Resume(ctx context.Context) error {
	q.Lock()
	q.paused = false
	q.Unlock()

	select {
	case q.ready <- struct{}{}:
	default:
	}
	return nil
}

func (q *queue) signal(ctx context.Context) error {
	q.Lock()
	count := len(q.workers)
	pause := q.paused
	q.Unlock()
	if pause {
		return nil
	}

	if count == 0 {
		return nil
	}

	builds, err := q.store.ListIncomplete(ctx)
	if err != nil {
		return err
	}

	q.Lock()
	defer q.Unlock()

	for _, build := range builds {
		if build.Status == core.StatusRunning {
			continue
		}
	loop:
		for w := range q.workers {
			//TODO filter

			select {
			case w.channel <- build:
				delete(q.workers, w)
				break loop
			}
		}
	}

	return nil
}
