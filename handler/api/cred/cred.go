package cred

import (
	"dubhe-ci/service/repos"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
)

func New(credService *repos.CredentialService) *CredentialHandler {
	return &CredentialHandler{credService: credService}
}

type CredentialHandler struct {
	credService *repos.CredentialService
}

func (c *CredentialHandler) Route(g *gin.RouterGroup) {
	g.POST("/cred", c.Create)
	g.DELETE("/cred/:id", c.Delete)
	g.GET("/cred", c.List)
	g.GET("/cred/rand", c.RandomGenerateSshKey)
}

func (c *CredentialHandler) Create(ctx *gin.Context) {
	var item *pb.Cred
	if err := utils.ParseJSON(ctx, &item); err != nil {
		utils.ErrorErr(ctx, err)
		return
	}

	_, err := c.credService.Create(utils.NewContext(ctx), item)
	if err != nil {
		utils.ErrorErr(ctx, err)
		return
	}

	utils.Success(ctx, nil)
}

func (c *CredentialHandler) Delete(ctx *gin.Context) {
	id := utils.S(ctx.Param("id")).String()

	_, err := c.credService.Delete(utils.NewContext(ctx), &pb.Id{Id: id})
	if err != nil {
		utils.ErrorErr(ctx, err)
		return
	}

	utils.Success(ctx, nil)
}

func (c *CredentialHandler) List(ctx *gin.Context) {

	records, err := c.credService.List(utils.NewContext(ctx), nil)
	if err != nil {
		utils.ErrorErr(ctx, err)
		return
	}

	utils.Success(ctx, records.Records)
}

func (c *CredentialHandler) RandomGenerateSshKey(ctx *gin.Context) {
	key, err := c.credService.RandomGenerateSshKey(utils.NewContext(ctx), nil)
	if err != nil {
		utils.ErrorErr(ctx, err)
		return
	}
	utils.Success(ctx, key)
}
