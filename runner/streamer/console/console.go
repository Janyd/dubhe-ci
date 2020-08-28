package console

import (
	"context"
	"dubhe-ci/core"
	"io"
	"os"
)

type Console struct {
	seq *Sequence
	col *Sequence
	tty bool
}

// New returns a new console recorder.
func New(tty bool) *Console {
	return &Console{
		tty: tty,
		seq: new(Sequence),
		col: new(Sequence),
	}
}

func (s *Console) Stream(_ context.Context, _ *core.State, name string) io.WriteCloser {
	if s.tty {
		return &pretty{
			base:  os.Stdout,
			color: colors[s.col.Next()%len(colors)],
			name:  name,
			seq:   s.seq,
		}
	}

	return &plain{
		base: os.Stdout,
		name: name,
		seq:  s.seq,
	}
}
