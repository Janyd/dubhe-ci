package socket

import (
	"context"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
)

const Broadcast = "Broadcast"
const LoadJob = "LoadJob"

func New() *Events {
	return &Events{}
}

type (
	EventRegister interface {
		Register(ctx context.Context, server *socketio.Server)
	}

	Events struct {
		Registers []EventRegister
	}
)

func (e *Events) Register(ctx context.Context, server *socketio.Server) {

	server.OnConnect("/", func(conn socketio.Conn) error {
		conn.SetContext("")
		logrus.WithField("id", conn.ID()).Infoln("connected!")

		conn.Join(Broadcast)
		return nil
	})
	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		logrus.WithField("reason", reason).Infoln("disconnected!")
		conn.Leave(Broadcast)
	})

	for _, register := range e.Registers {
		register.Register(ctx, server)
	}
}

func (e *Events) Add(register EventRegister) {
	e.Registers = append(e.Registers, register)
}
