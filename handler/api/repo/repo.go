package repo

import (
	"dubhe-ci/service/repos"
	"dubhe-ci/service/rpc/pb"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New(repoService *repos.RepositoryService) *RepositoryHandler {
	return &RepositoryHandler{repoService: repoService}
}

type RepositoryHandler struct {
	repoService *repos.RepositoryService
}

func (r *RepositoryHandler) Route(g *gin.RouterGroup) {

	g.POST("/repo", r.Create)
	g.GET("/repo/:id", r.Find)
	g.PUT("/repo", r.Update)
	g.DELETE("/repo/:id", r.Delete)
	g.POST("/repo/list", r.List)
	g.POST("/repo/scan/:repoId", r.Scan)
	g.POST("/repo/build/:repoId/:branchId", r.Build)

}

func (r *RepositoryHandler) Create(c *gin.Context) {
	var item *pb.Repo
	if err := utils.ParseJSON(c, &item); err != nil {
		utils.ErrorErr(c, err)
		return
	}

	repo, err := r.repoService.Create(utils.NewContext(c), item)
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, repo)
}

func (r *RepositoryHandler) Find(c *gin.Context) {

	id := utils.S(c.Param("id")).String()
	repo, err := r.repoService.Find(utils.NewContext(c), &pb.Id{Id: id})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, repo)
}

func (r *RepositoryHandler) Update(c *gin.Context) {
	var item *pb.Repo
	if err := utils.ParseJSON(c, &item); err != nil {
		utils.ErrorErr(c, err)
		return
	}

	_, err := r.repoService.Update(utils.NewContext(c), item)
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, nil)
}

func (r *RepositoryHandler) Delete(c *gin.Context) {
	id := utils.S(c.Param("id")).String()
	_, err := r.repoService.Delete(utils.NewContext(c), &pb.Id{Id: id})
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, nil)
}

func (r *RepositoryHandler) List(c *gin.Context) {
	var page *pb.Page
	if err := utils.ParseJSON(c, &page); err != nil {
		utils.ErrorErr(c, err)
		return
	}

	records, err := r.repoService.List(utils.NewContext(c), page)
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	utils.Success(c, records)
}

func (r *RepositoryHandler) Scan(c *gin.Context) {
	repoId := utils.S(c.Param("repoId")).String()
	err := r.repoService.Scan(utils.NewContext(c), repoId)
	if err != nil {
		logrus.WithError(err).WithField("repo.id", repoId).Errorln("scan fail")
	}
	utils.Success(c, nil)
}

func (r *RepositoryHandler) Build(c *gin.Context) {
	repoId := utils.S(c.Param("repoId")).String()
	branchId := utils.S(c.Param("branchId")).String()

	err := r.repoService.Build(utils.NewContext(c), repoId, branchId, true)
	if err != nil {
		logrus.WithError(err).WithField("repo.id", repoId).WithField("branch", branchId).Errorln("build fail")
	}

	utils.Success(c, nil)
}
