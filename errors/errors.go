package errors

// 定义通用错误
var (
	ErrCode = map[int]string{
		//正确编码
		200000: "操作成功",

		403000: "禁止访问",
		404000: "资源不存在",
		401000: "未授权",
		405000: "方法不支持",

		999999: "程序内部错误",
		900000: "无法找到内容",

		999901: "会话已过期，请重新登录",
		100001: "用户名密码错误",
		100002: "用户已被禁用",

		//参数错误
		100401: "参数错误",

		300000: "无实时日志",
		300001: "日志已完成",

		//业务代码
		500000: "系统繁忙，请重新尝试",
		500001: "名称已存在",
		500002: "用户名密码不可为空",
		500003: "公钥私钥不可为空",
	}
)

func ErrMsg(code int) string {
	msg, ok := ErrCode[code]
	if !ok {
		msg = ErrCode[999999]
	}
	return msg
}

// MessageError 自定义消息错误
type MessageError struct {
	err       error
	code      int
	msg       string
	secondMsg string
}

//集中判断常见错误
func Error(err error) *MessageError {
	return New(999999)
}

func New(code int) *MessageError {
	return NewMessageError(nil, code)
}

func NewMessageError(parent error, code int) *MessageError {
	return &MessageError{
		err:  parent,
		code: code,
		msg:  ErrMsg(code),
	}
}

func NewSecondMsg(parent error, code int, msg string) *MessageError {
	return &MessageError{
		err:       parent,
		code:      code,
		msg:       ErrMsg(code),
		secondMsg: msg,
	}
}

func (m *MessageError) Code() int {
	return m.code
}

func (m *MessageError) Error() string {
	return ErrCode[m.code]
}

func (m *MessageError) Parent() error {
	return m.err
}

func (m *MessageError) SecondMsg() string {
	return m.secondMsg
}
