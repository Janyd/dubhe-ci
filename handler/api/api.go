package api

import (
	"github.com/gin-gonic/gin"
)

func New() *Routers {
	return &Routers{
		Routers: make([]Router, 0),
	}
}

type (
	Router interface {
		Route(g *gin.RouterGroup)
	}

	Routers struct {
		Routers []Router
	}
)

func (r *Routers) Route(g *gin.RouterGroup) {
	for _, router := range r.Routers {
		router.Route(g)
	}
}

func (r *Routers) Add(router Router) {
	r.Routers = append(r.Routers, router)
}
