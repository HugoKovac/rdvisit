package errors

import (
	"errors"
	"fmt"
	"runtime"
)

type IErrorWrapper interface {
	FormatTrace() error
}

type ErrorWrapper struct {
	error
	file string
	line int
}

func (te *ErrorWrapper) FormatTrace() error {
	return fmt.Errorf("%s:%d - %w", te.file, te.line, te.error)
}

func Wrap(v any) error {
	if v == nil {
		return nil
	}
	_, file, line, ok := runtime.Caller(1)
	var err error
	switch e := v.(type) {
	case error:
		err = e
		if !ok {
			return &ErrorWrapper{error: err}
		}
	case string:
		err = errors.New(e)
	default:
		err = fmt.Errorf("%v", e)
	}
	return &ErrorWrapper{error: err, file: file, line: line}
}

func (te *ErrorWrapper) Unwrap() error {
	return te.error
}
