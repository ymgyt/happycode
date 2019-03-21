package errors

import "github.com/ymgyt/happycode/core/errors/errcode"

type Error struct {
	Code    errcode.Type
	Message string
	Inner   error
}

func (e *Error) Error() string {
	msg := e.Code.String()
	if e.Message != "" {
		msg += ": " + e.Message
	}
	if e.Inner != nil {
		msg += "\n" + e.Inner.Error()
	}
	return msg
}

func New(code errcode.Type) error {
	return &Error{Code: code}
}

func NewM(code errcode.Type, msg string) error {
	return &Error{Code: code, Message: msg}
}

func Wrap(code errcode.Type, err error) error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Inner: err}
}

func WrapM(code errcode.Type, err error, msg string) error {
	if err == nil {
		return nil
	}
	return &Error{Code: code, Inner: err, Message: msg}
}

func Code(err error) errcode.Type {
	if err == nil {
		return errcode.OK
	}
	myerr, ok := err.(*Error)
	if ok {
		return myerr.Code
	}
	return errcode.UnHandled
}
