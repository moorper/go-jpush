package jpush

import (
	"fmt"
)

// ErrorResponse 自定义错误
type ErrorResponse struct {
	ErrorMessage ErrorStruct `json:"error"`
}
type ErrorStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintln("code :", e.ErrorMessage.Code, " ", "message :", e.ErrorMessage.Message)
}
