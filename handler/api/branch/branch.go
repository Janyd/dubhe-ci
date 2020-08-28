package branch

import (
	"dubhe-ci/service/repos"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
)

func New(branchService *repos.BranchService) *BranchesHandler {
	return &BranchesHandler{branchService: branchService}
}

type BranchesHandler struct {
	branchService *repos.BranchService
}

func (b *BranchesHandler) Route(g *gin.RouterGroup) {
	g.GET("/branch/:repoId", b.List)
}

func (b *BranchesHandler) List(c *gin.Context) {

	repoId := utils.S(c.Param("repoId")).String()

	branches, err := b.branchService.List(utils.NewContext(c), &pb.Id{Id: repoId})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, branches.Records)
}
