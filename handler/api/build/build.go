package build

import (
	"dubhe-ci/service/repos"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
)

func New(buildService *repos.BuildService) *BuildsHandler {
	return &BuildsHandler{buildService: buildService}
}

type BuildsHandler struct {
	buildService *repos.BuildService
}

func (b *BuildsHandler) Route(g *gin.RouterGroup) {
	g.POST("/build/list/:repoId/:branchId", b.List)
	g.DELETE("/build/:buildId", b.Delete)
	g.GET("/build/steps/:buildId", b.ListStep)
}

func (b *BuildsHandler) List(c *gin.Context) {

	repoId := utils.S(c.Param("repoId")).String()
	branchId := utils.S(c.Param("branchId")).String()
	var page *pb.Page
	if err := utils.ParseJSON(c, &page); err != nil {
		utils.ErrorErr(c, err)
		return
	}

	records, err := b.buildService.List(utils.NewContext(c), &pb.BuildListRequest{
		Page:     page,
		RepoId:   repoId,
		BranchId: branchId,
	})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, records)
}

func (b *BuildsHandler) Delete(c *gin.Context) {
	buildId := utils.S(c.Param("buildId")).String()

	_, err := b.buildService.Delete(utils.NewContext(c), &pb.Id{Id: buildId})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, nil)
}

func (b *BuildsHandler) ListStep(c *gin.Context) {
	buildId := utils.S(c.Param("buildId")).String()

	steps, err := b.buildService.ListStep(utils.NewContext(c), &pb.Id{Id: buildId})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, steps.Records)
}
