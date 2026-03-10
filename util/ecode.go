package util

import "stock/constant"

// HttpCode 统一 HTTP 业务响应体：code + message + data，与前端约定一致
type HttpCode struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功返回：code=0，message 取常量表，data 可为 nil
func Success(data interface{}) HttpCode {
	return HttpCode{
		Code:    constant.CodeSuccess,
		Message: constant.GetMessage(constant.CodeSuccess),
		Data:    data,
	}
}

// Fail 失败返回：code 为业务错误码，message 从常量表取，data 为 nil
func Fail(code int) HttpCode {
	return HttpCode{
		Code:    code,
		Message: constant.GetMessage(code),
		Data:    nil,
	}
}

// FailWithData 失败但需要返回少量额外数据时使用
func FailWithData(code int, data interface{}) HttpCode {
	return HttpCode{
		Code:    code,
		Message: constant.GetMessage(code),
		Data:    data,
	}
}

// FailWithMessage 失败且需覆盖默认文案时使用（如参数校验提示）
func FailWithMessage(code int, message string) HttpCode {
	return HttpCode{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

// IsSuccess 判断是否为成功码
func (h HttpCode) IsSuccess() bool {
	return h.Code == constant.CodeSuccess
}

// EnsureMessage 若 Message 为空则按 Code 从常量表填充，便于兼容只设置 Code/Data 的旧写法
func (h HttpCode) EnsureMessage() HttpCode {
	if h.Message == "" {
		h.Message = constant.GetMessage(h.Code)
	}
	return h
}
