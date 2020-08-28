package utils

import (
	"dubhe-ci/errors"
	"github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

// 定义JSON操作
var (
	json              = jsoniter.ConfigCompatibleWithStandardLibrary
	JSONMarshal       = json.Marshal
	JSONUnmarshal     = json.Unmarshal
	JSONMarshalIndent = json.MarshalIndent
	JSONNewDecoder    = json.NewDecoder
	JSONNewEncoder    = json.NewEncoder
)

// JSONMarshalToString JSON编码为字符串
func JSONMarshalToString(v interface{}) string {
	s, err := jsoniter.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

func JsonResult(responseBody ResponseBody) string {
	buf, err := JSONMarshal(responseBody)
	if err != nil {
		logrus.WithError(err).WithField("object", responseBody).Error("响应JSON数据")
		return RErrorCode(999999)
	}
	return string(buf)

}

func RErrorCode(code int) string {
	responseBody := ResponseBody{
		Code: code,
		Msg:  errors.ErrMsg(code),
	}
	return JsonResult(responseBody)
}

func RError(err error) string {
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
	return JsonResult(responseBody)
}

func RSuccess(data interface{}) string {
	responseBody := ResponseBody{
		Code: SuccessCode,
		Data: data,
		Msg:  errors.ErrMsg(SuccessCode),
	}
	return JsonResult(responseBody)
}
