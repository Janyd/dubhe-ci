package runner

import (
	"context"
	"dubhe-ci/core"
	"io"
)

type Streamer interface {
	Stream(context.Context, *core.State, string) io.WriteCloser
}

func NopStreamer() Streamer {
	return new(nopStreamer)
}

type nopStreamer struct {
}

func (*nopStreamer) Stream(ctx context.Context, state *core.State, s string) io.WriteCloser {
	return new(nopWriteCloser)
}

type nopWriteCloser struct{}

func (*nopWriteCloser) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (*nopWriteCloser) Close() error {
	return nil
}
