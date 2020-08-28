package user

import (
	"dubhe-ci/config"
	"dubhe-ci/handler/auth"
	"dubhe-ci/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func New(config *config.Config, jwt *auth.JWT) *UsersHandler {
	return &UsersHandler{
		superUser: config.SuperUser,
		jwt:       jwt,
	}
}

type (
	UsersHandler struct {
		superUser config.SuperUser
		jwt       *auth.JWT
	}

	Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
)

func (u *UsersHandler) Route(g *gin.RouterGroup) {
	g.POST("/user/login", u.Login)
	g.GET("/user/info", u.Info)
	g.POST("/user/logout", u.Logout)
}

func (u *UsersHandler) Login(c *gin.Context) {

	var login Login
	if err := utils.ParseJSON(c, &login); err != nil {
		utils.ErrorErr(c, err)
		return
	}

	if !(u.superUser.Username == login.Username && u.superUser.Password == login.Password) {
		utils.Error(c, 100001)
		return
	}

	userId := 1

	claims := auth.Claims{
		UserId: userId,
		RoleId: 0,
		StandardClaims: jwt.StandardClaims{
			Issuer:  strconv.Itoa(userId),
			Subject: strconv.Itoa(userId),
		},
	}
	token, err := u.jwt.CreateToken(claims)
	if err != nil {
		utils.ErrorErr(c, err)
		return
	}
	utils.Success(c, token)
}
func (u *UsersHandler) Info(c *gin.Context) {
	//TODO

	info := map[string]interface{}{
		"name":         "admin",
		"avatar":       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		"introduction": "I am a super administrator",
		"roles":        []string{"admin"},
	}
	utils.Success(c, info)
}

func (u *UsersHandler) Logout(c *gin.Context) {
	utils.Success(c, nil)
}
