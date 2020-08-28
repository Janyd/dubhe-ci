package middleware

import (
	"dubhe-ci/handler/auth"
	"dubhe-ci/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(j *auth.JWT, skipper ...SkipperFunc) gin.HandlerFunc {

	return func(c *gin.Context) {
		var userId int
		if t := utils.GetToken(c); t != "" {
			claims, err := j.ParseToken(t)
			if err != nil {
				utils.Error(c, 999901)
				return
			}
			userId = claims.UserId
		}

		if userId > 0 {
			c.Set(utils.UserIDKey, userId)
		}

		if len(skipper) > 0 && skipper[0](c) {
			c.Next()
			return
		}

		if userId == 0 {
			utils.Error(c, 999901)
		}
	}
}
