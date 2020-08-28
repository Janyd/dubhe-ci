package core

import (
	"sync"
	"time"
)

type State struct {
	sync.Mutex

	Build *Build
	Repo  *Repository
}

func (s *State) Cancel() {
	s.Lock()
	s.skipAll()
	s.killAll()
	s.update()
	s.Unlock()
}

func (s *State) Cancelled() bool {
	s.Lock()
	v := s.killed()
	s.Unlock()
	return v
}

func (s *State) skipAll() {
	for _, v := range s.Build.Steps {
		s.skip(v)
	}
}

func (s *State) skip(v *Step) {
	if v.Status == StatusPending {
		v.Started = time.Now()
		v.Stopped = time.Now()
		v.Status = StatusSkipped
		v.ExitCode = 0
		v.Error = ""
	}
}

func (s *State) Skipped() bool {
	s.Lock()
	v := s.skipped()
	s.Unlock()
	return v
}

func (s *State) Start(name string) {
	s.Lock()
	v := s.Find(name)
	s.start(v)
	s.Unlock()
}

func (s *State) Finish(name string, code int) {
	s.Lock()
	v := s.Find(name)
	s.finish(v, code)
	s.update()
	s.Unlock()
}

func (s *State) finish(v *Step, code int) {
	switch v.Status {
	case StatusRunning, StatusPending:
	default:
		return
	}
	v.ExitCode = code
	v.Stopped = time.Now()
	switch code {
	case 0, 78:
		v.Status = StatusPassing
	default:
		v.Status = StatusFailing
	}
}

func (s *State) finished() bool {
	for _, v := range s.Build.Steps {
		switch v.Status {
		case StatusRunning, StatusPending:
			return false
		}
	}
	return true
}

func (s *State) skipped() bool {
	if s.finished() == false {
		return false
	}

	for _, v := range s.Build.Steps {
		if v.Status == StatusSkipped {
			return true
		}
	}
	return false
}

func (s *State) kill(v *Step) {
	if v.Status == StatusRunning {
		v.Status = StatusKilled
		v.Stopped = time.Now()
		v.ExitCode = 137
		v.Error = ""
	}
}

func (s *State) killAll() {
	s.Build.Error = ""
	//s.Build.ExitCode = 0
	s.Build.Status = StatusKilled
	s.Build.Finished = time.Now()

	for _, v := range s.Build.Steps {
		s.skip(v)
	}
}

func (s *State) killed() bool {
	return s.Build.Status == StatusKilled
}

func (s *State) start(v *Step) {
	if v.Status == StatusPending {
		v.Status = StatusRunning
		v.Started = time.Now()
		v.ExitCode = 0
		v.Error = ""
	}
}

func (s *State) failAll(err error) {
	switch s.Build.Status {
	case StatusPending, StatusRunning:
		s.Build.Status = StatusError
		s.Build.Error = err.Error()
		s.Build.Finished = time.Now()
		//s.Build.ExitCode = 255
	}
}

func (s *State) SkipAll() {
	s.Lock()
	s.skipAll()
	s.update()
	s.Unlock()
}

func (s *State) Skip(name string) {
	s.Lock()
	v := s.Find(name)
	s.skip(v)
	s.update()
	s.Unlock()
}

func (s *State) Failed() bool {
	s.Lock()
	v := s.failed()
	s.Unlock()
	return v
}

func (s *State) failed() bool {
	switch s.Build.Status {
	case StatusFailing,
		StatusError,
		StatusKilled:
		return true
	}
	for _, v := range s.Build.Steps {
		//if v.ErrIgnore {
		//	continue
		//}
		switch v.Status {
		case StatusFailing,
			StatusError,
			StatusKilled:
			return true
		}
	}
	return false
}

// FailAll fails the entire pipeline.
func (s *State) FailAll(err error) {
	s.Lock()
	s.failAll(err)
	s.skipAll()
	s.update()
	s.Unlock()
}

func (s *State) update() {
	for _, v := range s.Build.Steps {
		switch v.Status {
		case StatusKilled:
			//s.Build.ExitCode = 137
			s.Build.Status = StatusKilled
			s.Build.Status = StatusKilled
			return
		case StatusError:
			s.Build.Error = v.Error
			//s.Build.ExitCode = 255
			s.Build.Status = StatusError
			s.Build.Status = StatusError
			return
			//case core.StatusFailing:
			//	if v.ErrIgnore == false {
			//		s.Build.Status = core.StatusFailing
			//		return
			//	}
		}
	}
}

func (s *State) Find(name string) *Step {
	for _, step := range s.Build.Steps {
		if step.Name == name {
			return step
		}
	}
	panic("step not found: " + name)
}

func (s *State) FinishAll() {
	s.Lock()
	s.finishAll()
	s.update()
	s.Unlock()
}

func (s *State) finishAll() {
	for _, v := range s.Build.Steps {
		switch v.Status {
		case StatusPending:
			s.skip(v)
		case StatusRunning:
			s.finish(v, 0)
		}
	}
	switch s.Build.Status {
	case StatusRunning, StatusPending:
		s.Build.Finished = time.Now()
		s.Build.Status = StatusPassing
		if s.failed() {
			s.Build.Status = StatusFailing
		}
	default:
		s.Build.Finished = time.Now()
	}
}
