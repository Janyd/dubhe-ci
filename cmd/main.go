package main

import (
	"context"
	"dubhe-ci/config"
	"dubhe-ci/logger"
	"dubhe-ci/runner"
	"dubhe-ci/server"
	"dubhe-ci/utils"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "配置文件(.json,.yaml,.toml)")
	flag.StringVar(&configFile, "c", "", "配置文件(.json,.yaml,.toml)")
}

func main() {
	flag.Parse()

	if configFile == "" {
		logrus.Fatalln("请使用-c或-config指定配置文件")
	}

	c, err := parseConfig(configFile)
	if err != nil {
		logrus.Fatalln("解析配置文件失败", err)
	}

	err = logger.InitLogger(c.Log)
	if err != nil {
		logrus.Fatalln(err)
	}

	err = utils.InitSnowflake()
	if err != nil {
		logrus.Fatalln(err)
	}

	ctx := withContextCancel(context.Background())

	app, err := InitializeApplication(c)
	if err != nil {
		logrus.WithError(err).Fatalln("cannot initialize grpcServer")
	}

	g := errgroup.Group{}

	g.Go(func() error {
		return app.socketServer.ListenAndServe(ctx)
	})

	g.Go(func() error {
		return app.httpServer.ListenAndServe(ctx)
	})

	g.Go(func() error {
		return app.grpcServer.ListenAndServe(ctx)
	})

	g.Go(func() error {
		if app.runner == nil {
			return nil
		}

		logrus.WithField("threads", c.Capacity).
			Infoln("main: starting the local build runner")
		return app.runner.Start(ctx, c.Capacity)
	})

	if err := g.Wait(); err != nil {
		logrus.WithError(err).Fatalln("program terminated")
	}
}

func withContextCancel(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			println("interrupt received, terminating process")
			cancel()
		}
	}()

	return ctx
}

func parseConfig(filepath string) (*config.Config, error) {
	var c config.Config
	_, err := toml.DecodeFile(filepath, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type application struct {
	socketServer *server.SocketIOServer
	grpcServer   *server.GrpcServer
	httpServer   *server.GinServer
	runner       *runner.Runner
}

func newApplication(
	socketServer *server.SocketIOServer,
	grpcServer *server.GrpcServer,
	httpServer *server.GinServer,
	runner *runner.Runner,

) application {
	return application{
		socketServer: socketServer,
		httpServer:   httpServer,
		grpcServer:   grpcServer,
		runner:       runner,
	}
}
