package server

import (
	"context"
	"dubhe-ci/socket"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"golang.org/x/sync/errgroup"
)

type SocketIOServer struct {
	Server *socketio.Server
	Events *socket.Events
}

func (s SocketIOServer) Route(g *gin.RouterGroup) {
	g.GET("/socket.io/", func(c *gin.Context) {
		s.Server.ServeHTTP(c.Writer, c.Request)
	})
}

func (s SocketIOServer) ListenAndServe(ctx context.Context) error {
	s.Events.Register(ctx, s.Server)

	var g errgroup.Group
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return s.Server.Close()
		}
	})

	g.Go(func() error {
		return s.Server.Serve()
	})

	return g.Wait()
}
