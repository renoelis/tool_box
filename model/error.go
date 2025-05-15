package model

// 错误响应
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 实现error接口
func (e *ErrorResponse) Error() string {
	return e.Message
} 