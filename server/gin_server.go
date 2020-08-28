package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type GinServer struct {
	Addr    string
	Handler http.Handler
}

func (s GinServer) ListenAndServe(ctx context.Context) error {
	var g errgroup.Group
	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return srv.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		logrus.WithField("addr", s.Addr).Infoln("starting the http server")
		return srv.ListenAndServe()
	})

	return g.Wait()
}
