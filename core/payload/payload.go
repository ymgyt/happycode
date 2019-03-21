package payload

import (
	"bytes"
	"encoding/gob"
)

func init() {
	gob.Register(Hello{})
}

type Type int

const (
	TypeHello Type = iota
)

func (t Type) String() string {
	switch t {
	case TypeHello:
		return "hello"
	}
	return "undefined"
}

type Interface interface {
	Type() Type
}

type Hello struct {
	Message string
}

func (h Hello) Type() Type { return TypeHello }

func Encode(p Interface) []byte {
	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(&p)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}
