package runner

import (
	"dubhe-ci/runner/compiler"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
)

type State struct {
	hook *Hook

	spec  *compiler.Spec
	step  *compiler.Step
	state *compiler.State
}

func snapshot(e *Execer, step *compiler.Step, state *compiler.State) *State {
	s := &State{
		hook:  e.hook,
		spec:  e.spec,
		step:  step,
		state: state,
	}
	return s
}

func (s *State) output(rc io.ReadCloser) error {
	w := NewWriter(s)
	stdcopy.StdCopy(w, w, rc)

	if w.state.hook.GotLogs != nil {
		return w.state.hook.GotLogs(s, w.lines)
	}

	return rc.Close()
}
