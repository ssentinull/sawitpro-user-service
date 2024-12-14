package utils

import (
	"fmt"
	"math"
	"net/http"
	"runtime"
)

type ErrorCode uint32

type errorMessage map[ErrorCode]string

var (
	errorMessages = errorMessage{
		http.StatusBadRequest:          `Invalid Input. Please Validate Your Input.`,
		http.StatusUnauthorized:        `Unauthorized Access. You are not authorized to access this resource.`,
		http.StatusForbidden:           `Forbidden Access. You are forbidden to access this resource`,
		http.StatusNotFound:            `Record Does Not Exist. Please Validate Your Input Or Contact Administrator.`,
		http.StatusConflict:            `Record Has Existed and Must Be Unique. Please Validate Your Input Or Contact Administrator.`,
		http.StatusUnprocessableEntity: `Unprocessable Entity. This entity can not be processed.`,
		http.StatusInternalServerError: `Internal Server Error. Please Call Administrator.`,
	}
)

const NoCode ErrorCode = math.MaxUint32

type stacktrace struct {
	message string
	cause   error
	code    ErrorCode
	file    string
	line    int
}

func Wrap(cause error, msg string, vals ...interface{}) error {
	if cause == nil {
		return nil
	}
	return create(cause, NoCode, msg, vals...)
}

func New(msg string, vals ...interface{}) error {
	return create(nil, NoCode, msg, vals...)
}

func NewErrorWithCode(code ErrorCode, msg string, vals ...interface{}) error {
	return create(nil, code, msg, vals...)
}

func WrapWithCode(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if cause == nil {
		return nil
	}
	return create(cause, code, msg, vals...)
}

func create(cause error, code ErrorCode, msg string, vals ...interface{}) error {
	if code == 0 {
		code = http.StatusInternalServerError
	}

	err := &stacktrace{
		message: fmt.Sprintf(msg, vals...),
		cause:   cause,
		code:    code,
	}

	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return err
	}

	err.file, err.line = file, line
	f := runtime.FuncForPC(pc)
	if f == nil {
		return err
	}

	return err
}

func GetCode(err error) ErrorCode {
	if err, ok := err.(*stacktrace); ok {
		return err.code
	}
	return NoCode
}

func GetCause(err error) error {
	if err, ok := err.(*stacktrace); ok {
		return err.cause
	}
	return err
}

func GetMessage(err error) string {
	errStacktrace, ok := err.(*stacktrace)
	if !ok {
		return ""
	}

	if errStacktrace.message == "" {
		return errorMessages[errStacktrace.code]
	}

	return errStacktrace.message
}

func (st *stacktrace) Error() string {
	if st != nil {
		return fmt.Sprint(st)
	}

	return ""
}
