package errorx

const successCode = 0
const successMsg = "Success"
const defaultCode = 1001
const paramsFailedCode = 1002
const categoryError = 1003
const NotFound = 1004

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// CodeErrorResponse 此处结构体用于统一返回结果
type CodeErrorResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewCategoryError(msg string) error {
	return NewCodeError(categoryError, msg)
}

func NewCodeError(code int, msg string) error {
	return &CodeError{Code: code, Msg: msg}
}

// NewDefaultError 设置默认错误码
func NewDefaultError(msg string) error {
	return NewCodeError(defaultCode, msg)
}

// NewParamsFailedError 设置接口参数json校验错误码
func NewParamsFailedError(msg string) error {
	return NewCodeError(paramsFailedCode, msg)
}

func NewDefaultParamsFailedError() error {
	return NewCodeError(paramsFailedCode, "Params parse failed.")
}

func (e *CodeError) Error() string {
	return e.Msg
}

// Data /**返回通用错误数据结构
func (e *CodeError) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: e.Code,
		Msg:  e.Msg,
	}
}

// NewSuccessJson /**接口请求成功返回数据
func NewSuccessJson(resp interface{}) *CodeErrorResponse {
	return &CodeErrorResponse{
		Code: successCode,
		Msg:  successMsg,
		Data: resp,
	}
}
