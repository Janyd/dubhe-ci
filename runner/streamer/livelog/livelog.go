package livelog

import (
	"bytes"
	"context"
	"dubhe-ci/core"
	"errors"
	"github.com/sirupsen/logrus"
	"sync"
)

var errStreamNotFound = errors.New("stream: not found")

func New(logs core.LogStore, steps core.StepStore) core.LogStream {
	return &streamer{
		streams: make(map[string]*stream),
		logs:    logs,
		steps:   steps,
	}
}

type (
	streamer struct {
		sync.Mutex

		streams map[string]*stream
		logs    core.LogStore
		steps   core.StepStore
	}
)

func (s *streamer) Create(ctx context.Context, stepId string) error {
	s.Lock()
	s.streams[stepId] = newStream()
	s.Unlock()
	return nil
}

func (s *streamer) Delete(ctx context.Context, stepId string) error {
	s.Lock()
	stream, ok := s.streams[stepId]
	if ok {
		delete(s.streams, stepId)
	}

	raw := stream.all()
	s.Unlock()
	if !ok {
		return errStreamNotFound
	}

	step, err := s.steps.Find(ctx, stepId)
	if err != nil {
		logrus.WithField("step.id", stepId).WithError(err).
			Errorln("cannot find step from db")
		return err
	}
	buf := bytes.NewBuffer(raw)
	err = s.logs.Create(ctx, step.BuildId, step.Id, buf)
	if err != nil {
		return err
	}

	return stream.close()
}

func (s *streamer) Write(ctx context.Context, stepId string, line *core.Line) error {
	s.Lock()
	stream, ok := s.streams[stepId]
	s.Unlock()
	if !ok {
		return errStreamNotFound
	}

	return stream.write(line)
}

func (s *streamer) Tail(ctx context.Context, stepId string, sub core.Subscriber) bool {
	s.Lock()
	stream, ok := s.streams[stepId]
	s.Unlock()
	if !ok {
		return false
	}
	stream.subscribe(ctx, sub)
	return true
}
