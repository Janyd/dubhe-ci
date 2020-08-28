package socket

import (
	"context"
	"dubhe-ci/core"
	"dubhe-ci/utils"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

func NewStepEvent(stream core.LogStream) *StepEvent {
	return &StepEvent{
		stream: stream,
	}
}

type (
	StepEvent struct {
		ctx    context.Context
		stream core.LogStream
	}
	StepJoin struct {
		StepId string `json:"stepId"`
	}

	StepSubscriber struct {
		conn   socketio.Conn
		stepId string
	}
)

func (s *StepSubscriber) Publish(line *core.Line) {
	s.conn.Emit("step:"+s.stepId, utils.RSuccess(line))
}

func (s *StepSubscriber) Close() {
	s.conn.Emit("step:"+s.stepId, utils.RErrorCode(300001))
}

func (s *StepEvent) Register(ctx context.Context, server *socketio.Server) {
	server.OnEvent("", "step:join", func(conn socketio.Conn, data string) string {
		log := logrus.WithField("clientId", conn.ID())
		var p StepJoin

		err := utils.JSONUnmarshal([]byte(data), &p)

		if err != nil {
			log.Warnln("parse json error")
			return utils.RError(err)
		}
		log = log.WithField("stepId", p.StepId)

		sub := &StepSubscriber{
			conn:   conn,
			stepId: p.StepId,
		}

		if s.stream.Tail(ctx, p.StepId, sub) {
			return utils.RSuccess(nil)
		}

		return utils.RErrorCode(300000)
	})
}
