package runner

import (
	"dubhe-ci/core"
	"strings"
	"time"
)

type lineWriter struct {
	num   int
	now   time.Time
	lines []*core.Line
	size  int
	limit int

	state *State
}

func NewWriter(state *State) *lineWriter {
	return &lineWriter{
		num:   0,
		now:   time.Now(),
		limit: 5242880,
		state: state,
	}
}

func (w *lineWriter) Write(p []byte) (n int, err error) {
	if w.size >= w.limit {
		return len(p), nil
	}
	s := string(p)
	s = strings.TrimSuffix(s, "\n")
	parts := strings.Split(s, "\n")

	for _, part := range parts {
		line := &core.Line{
			Number:    w.num,
			Message:   part,
			Timestamp: int64(time.Since(w.now).Seconds()),
		}

		if w.state.hook.GotLine != nil {
			_ = w.state.hook.GotLine(w.state, line)
		}
		w.size = w.size + len(part)
		w.num++

		w.lines = append(w.lines, line)
	}

	if w.size >= w.limit {
		w.lines = append(w.lines, &core.Line{
			Number:    w.num,
			Message:   "warning: maximum output exceeded",
			Timestamp: int64(time.Since(w.now).Seconds()),
		})
	}

	return len(p), nil
}
