package runner

import "dubhe-ci/core"

type Hook struct {
	//执行步骤前操作
	BeforeEach func(state *State) error

	//执行步骤后操作
	AfterEach func(state *State) error

	GotLine func(state *State, line *core.Line) error

	GotLogs func(state *State, lines []*core.Line) error
}
