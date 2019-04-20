package errors

import (
	"fmt"
	"strconv"
)

// Error is a generic error struct
type Error struct {
	Code        string `json:"error"`
	HTTPCode    int    `json:"-"`
	Description string `json:"error_description"`
}

func (err Error) Error() string {
	return fmt.Sprintf("code:%s,descprition:%s", err.Code, err.Description)
}

// New create a new error with a string err code
func New(code string, httpCode int, desc string) Error {
	return Error{
		Code:        code,
		HTTPCode:    httpCode,
		Description: desc,
	}
}

// NewI create a new error with a int err code
func NewI(code, httpCode int, desc string) Error {
	return New(strconv.Itoa(code), httpCode, desc)
}

// AppendDesc append additional description
func (err *Error) AppendDesc(msg string) {
	err.Description += ": " + msg
}
