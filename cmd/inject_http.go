package main

import (
	"dubhe-ci/config"
	"dubhe-ci/handler/api"
	"dubhe-ci/handler/api/branch"
	"dubhe-ci/handler/api/build"
	"dubhe-ci/handler/api/cred"
	"dubhe-ci/handler/api/logs"
	"dubhe-ci/handler/api/repo"
	"dubhe-ci/handler/api/user"
	"dubhe-ci/handler/auth"
	"dubhe-ci/handler/middleware"
	"dubhe-ci/server"
	"dubhe-ci/socket"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	engineio "github.com/googollee/go-engine.io"
	"github.com/googollee/go-engine.io/transport"
	"github.com/googollee/go-engine.io/transport/polling"
	"github.com/googollee/go-engine.io/transport/websocket"
	socketio "github.com/googollee/go-socket.io"
	"github.com/sirupsen/logrus"
	"net/http"
)

var ginSet = wire.NewSet(
	auth.NewJWT,
	user.New,
	repo.New,
	cred.New,
	branch.New,
	build.New,
	logs.New,
	socket.NewStepEvent,
	provideEvents,
	provideSocketIOServer,
	provideRouters,
	provideGinServer,
)

func provideRouters(
	repoCtl *repo.RepositoryHandler,
	credCtl *cred.CredentialHandler,
	branchCtl *branch.BranchesHandler,
	userHandler *user.UsersHandler,
	buildsHandler *build.BuildsHandler,
	logHandler *logs.LogHandler,
	socketServer *server.SocketIOServer,
) *api.Routers {
	routers := api.New()

	routers.Add(repoCtl)
	routers.Add(credCtl)
	routers.Add(branchCtl)
	routers.Add(userHandler)
	routers.Add(buildsHandler)
	routers.Add(logHandler)
	routers.Add(socketServer)

	return routers
}

func provideGinServer(config *config.Config, routers *api.Routers, jwt *auth.JWT) *server.GinServer {
	engine := gin.New()

	auth.SetSignKey(config.JWTAuth.SigningKey)
	auth.SetExpired(config.JWTAuth.Expired)

	engine.NoMethod(middleware.NoMethodHandler())
	engine.NoRoute(middleware.NoRouteHandler())

	// 崩溃恢复
	engine.Use(middleware.RecoveryMiddleware())
	//跨域
	engine.Use(middleware.CORSMiddleware())

	group := engine.Group("/api")
	group.Use(middleware.AuthMiddleware(
		jwt,
		middleware.AllowMethodAndPathPrefixSkipper(
			middleware.JoinRouter("GET", "/api/user/login"),
			middleware.JoinRouter("POST", "/api/user/login"),
			middleware.JoinRouter("GET", "/api/socket.io"),
			middleware.JoinRouter("POST", "/api/socket.io"),
		),
	))
	routers.Route(group)

	ginpprof.Wrap(engine)

	return &server.GinServer{
		Addr:    config.Http.Address,
		Handler: engine,
	}
}

func provideSocketIOServer(events *socket.Events) *server.SocketIOServer {

	wt := websocket.Default
	wt.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	s, err := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			polling.Default,
			wt,
		},
	})
	if err != nil {
		logrus.WithError(err).Fatalln("cannot create socketio server")
	}

	socketServer := &server.SocketIOServer{
		Server: s,
		Events: events,
	}

	return socketServer
}

func provideEvents(event *socket.StepEvent) *socket.Events {
	events := socket.New()
	events.Add(event)

	return events
}
