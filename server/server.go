package server

import (
	"context"
	"dubhe-ci/service/rpc"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type GrpcServer struct {
	Addr       string
	Network    string
	GrpcServer *grpc.Server
	Registers  *rpc.Registers
}

func (s GrpcServer) ListenAndServe(ctx context.Context) error {
	var g errgroup.Group
	logrus.Infoln("Application Starting...")
	s.GrpcServer = grpc.NewServer()
	reflection.Register(s.GrpcServer)
	g.Go(func() error {
		select {
		case <-ctx.Done():
			s.GrpcServer.Stop()
			return nil
		}
	})

	g.Go(func() error {
		s.Registers.Register(s.GrpcServer)
		listener, err := net.Listen(s.Network, s.Addr)
		if err != nil {
			return err
		}
		logrus.WithFields(logrus.Fields{
			"addr":  s.Addr,
			"proto": s.Network,
		}).Infoln("starting the rpc server")
		return s.GrpcServer.Serve(listener)
	})

	return g.Wait()
}
