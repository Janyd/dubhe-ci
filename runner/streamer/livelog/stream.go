package livelog

import (
	"context"
	"dubhe-ci/core"
	"encoding/json"
	"sync"
)

const bufferSize = 5000

type stream struct {
	sync.Mutex

	hist []*core.Line
	list map[core.Subscriber]struct{}
}

func newStream() *stream {
	return &stream{
		list: map[core.Subscriber]struct{}{},
	}
}

func (s *stream) write(line *core.Line) error {
	s.Lock()
	s.hist = append(s.hist, line)

	for l := range s.list {
		l.Publish(line)
	}

	if size := len(s.hist); size >= bufferSize {
		s.hist = s.hist[size-bufferSize:]
	}
	s.Unlock()
	return nil
}

func (s *stream) subscribe(ctx context.Context, sub core.Subscriber) {
	s.Lock()
	for _, line := range s.hist {
		sub.Publish(line)
	}

	s.list[sub] = struct{}{}
	s.Unlock()

	return
}
func (s *stream) close() error {
	s.Lock()
	defer s.Unlock()
	for sub := range s.list {
		delete(s.list, sub)
		sub.Close()
	}

	return nil
}

func (s *stream) all() []byte {
	s.Lock()
	defer s.Unlock()
	raw, _ := json.Marshal(s.hist)
	return raw
}
