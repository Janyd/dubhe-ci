package logs

import (
	"bytes"
	"dubhe-ci/core"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
)

func New(logsStore core.LogStore) *LogHandler {
	return &LogHandler{logsStore: logsStore}
}

type LogHandler struct {
	logsStore core.LogStore
}

func (l *LogHandler) Route(g *gin.RouterGroup) {
	g.GET("/log/:stepId", l.Find)
}

func (l *LogHandler) Find(c *gin.Context) {
	stepId := utils.S(c.Param("stepId")).String()
	rc, err := l.logsStore.Find(utils.NewContext(c), stepId)
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(rc)
	s := buf.String()

	utils.Success(c, s)
}
