package requester

import (
	"fmt"
)

const (
	CodeUnknown         int = 0
	CodeInvalidBody     int = 1
	CodeBadResponseCode int = 2
	CodeEncodingError   int = 3
	CodeInvalidForm     int = 4
	CodeMissingURL      int = 5
)

type statusError struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e statusError) Error() string {
	return fmt.Sprintf("request error: code = %v message = %v", e.Code, e.Message)
}

func Code(err error) int {
	if err == nil {
		return CodeUnknown
	}
	if e, ok := err.(statusError); ok {
		return e.Code
	}
	return CodeUnknown
}

func StatusCode(err error) int {
	if err == nil {
		return CodeUnknown
	}
	if e, ok := err.(statusError); ok {
		return e.StatusCode
	}
	return CodeUnknown
}
