package console

import (
	"fmt"
	"io"
	"strings"
)

const prettyf = "\033[%s[%s:%d]\033[0m %s\n"

var colors = []string{
	"32m", //green
	"33m", //yellow
	"34m", //blue
	"35m", //magenta
	"36m", //cyan
}

type pretty struct {
	base  io.Writer
	color string
	name  string
	seq   *Sequence
}

func (w *pretty) Write(p []byte) (n int, err error) {
	for _, part := range split(p) {
		_, _ = fmt.Fprintf(w.base, prettyf, w.color, w.name, w.seq.Next(), part)
	}
	return len(p), nil
}

func (w *pretty) Close() error {
	return nil
}

func split(b []byte) []string {
	s := string(b)
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
