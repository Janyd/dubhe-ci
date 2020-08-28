package utils

import (
	"context"
	"dubhe-ci/errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// 定义上下文中的键
const (
	prefix = "ginadmin"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user_id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace_id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res_body"

	SuccessCode = 200000
)

type ResponseBody struct {
	Code      int         `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
	SecondMsg string      `json:"secondMsg"`
}

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		logrus.WithError(err).Error("解析失败")
		return err
	}
	return nil
}

// NewContext get context.Context
func NewContext(c *gin.Context) context.Context {
	parent := context.Background()
	//TODO
	return parent
}

// GetToken 获取用户令牌
func GetToken(c *gin.Context) string {
	var token string
	auth := c.GetHeader("Authorization")
	prefix := "Bearer "
	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	}
	return token
}

func Success(c *gin.Context, data interface{}) {
	responseBody := ResponseBody{
		Code: SuccessCode,
		Data: data,
		Msg:  errors.ErrMsg(SuccessCode),
	}
	Result(c, http.StatusOK, responseBody)
}

func Error(c *gin.Context, code int) {
	ErrorStatus(c, http.StatusOK, code)
}

func ErrorErr(c *gin.Context, err error) {
	var msg string
	var code int
	var secondMsg string
	switch e := err.(type) {
	case *errors.MessageError:
		msg = e.Error()
		code = e.Code()
		secondMsg = e.SecondMsg()
	default:
		logrus.WithError(err).Error("未知响应错误")
		code = 999999
		msg = errors.ErrMsg(code)
	}
	responseBody := ResponseBody{
		Code:      code,
		Msg:       msg,
		SecondMsg: secondMsg,
	}
	Result(c, http.StatusOK, responseBody)
}

func ErrorStatus(c *gin.Context, status int, code int) {
	responseBody := ResponseBody{
		Code: code,
		Msg:  errors.ErrMsg(code),
	}
	Result(c, status, responseBody)
}

func Result(c *gin.Context, status int, responseBody ResponseBody) {
	buf, err := JSONMarshal(responseBody)
	if err != nil {
		logrus.WithError(err).WithField("object", responseBody).Error("响应JSON数据")
		Error(c, 999999)
		return
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}
