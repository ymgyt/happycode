package js

import (
	"errors"
	"fmt"
	"net/url"
	"syscall/js"
)

var (
	Null      = js.Null()
	Undefiled = js.Undefined()
	ZeroValue = js.Value{}

	ErrElementNotFound = errors.New("element not found")
)

type Type int

const (
	TypeObject Type = iota
)

func (typ Type) String() string {
	switch typ {
	case TypeObject:
		return "object"
	default:
		return "myundefined"
	}
}

func NewDocument() *Element {
	d := js.Global().Get("document")
	if isNull(d) {
		panic(`js.Global().Get("document") return null`)
	}
	return NewElement(d)
}

func BackendURL() *url.URL {
	// http://localhost:8888/
	href := js.Global().Get("location").Get("href").String()
	ep, err := url.Parse(href)
	if err != nil {
		panic(err)
	}
	return ep
}

func isNull(v js.Value) bool {
	return v == Null
}

func assertType(v js.Value, typ Type) {
	if v.Type().String() != typ.String() {
		panic(fmt.Errorf("assertType fail: got %s, want %s", v.Type().String(), typ.String()))
	}
}
