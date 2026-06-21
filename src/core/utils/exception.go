package utils

import (
	"errors"
	"github.com/golibs-starter/golib/exception"
)

func MakeException(code uint, errOrMsg interface{}) *exception.Exception {
	var err error

	switch v := errOrMsg.(type) {
	case error:
		err = v
	case string:
		err = errors.New(v)
	default:
		err = errors.New("unknown error")
	}

	ex := exception.New(code, err.Error())
	return &ex
}
