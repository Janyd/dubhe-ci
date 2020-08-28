package console

import (
	"fmt"
	"io"
)

const plainf = "[%s:%d] %s\n"

type plain struct {
	base io.Writer
	name string
	seq  *Sequence
}

func (w *plain) Write(b []byte) (n int, err error) {
	for _, part := range split(b) {
		_, _ = fmt.Fprintf(w.base, plainf, w.name, w.seq.Next(), part)
	}
	return len(b), nil
}

func (w *plain) Close() error {
	return nil
}
