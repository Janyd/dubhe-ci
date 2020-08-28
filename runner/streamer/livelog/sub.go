package livelog

import (
	"dubhe-ci/core"
	"sync"
)

type subscriber struct {
	sync.Mutex

	handler chan *core.Line
	closec  chan struct{}
	closed  bool
}

func (s *subscriber) publish(line *core.Line) {
	select {
	case <-s.closec:
	case s.handler <- line:
	default:
	}
}

func (s *subscriber) close() {
	s.Lock()
	if !s.closed {
		close(s.closec)
		s.closed = true
	}
	s.Unlock()
}
