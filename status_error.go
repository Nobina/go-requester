package requester

import (
	"fmt"
)

const (
	CodeUnknown           int32 = 0
	CodeInvalidBody       int32 = 1
	CodeBadResponseStatus int32 = 2
	CodeEncodingError     int32 = 3
	CodeInvalidForm       int32 = 4
	CodeMissingURL        int32 = 5
)

type statusError struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (e statusError) Error() string {
	return fmt.Sprintf("request error: code = %v message = %v", e.Code, e.Message)
}

func Code(err error) int32 {
	if err == nil {
		return CodeUnknown
	}
	if e, ok := err.(statusError); ok {
		return e.Code
	}
	return CodeUnknown
}
