package console

import "sync"

type Sequence struct {
	sync.Mutex
	value int
}

func (s *Sequence) Next() int {
	s.Lock()
	s.value++
	i := s.value
	s.Unlock()
	return i
}

func (s *Sequence) Curr() int {
	s.Lock()
	i := s.value
	s.Unlock()
	return i
}
