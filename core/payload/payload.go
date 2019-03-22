package payload

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/ymgyt/happycode/core/config"
)

func init() {
	gob.Register(Hello{})
	gob.Register(ConfigRequest{})
	gob.Register(ConfigResponse{})
}

type Type int

const (
	TypeHello Type = iota
	TypeConfigRequest
	TypeConfigResponse
)

func (t Type) String() string {
	switch t {
	case TypeHello:
		return "Hello"
	case TypeConfigRequest:
		return "ConfigRequest"
	case TypeConfigResponse:
		return "ConfigResponse"
	}
	return "Undefined"
}

type Interface interface {
	Type() Type
}

type Hello struct {
	Message string
}

func (h Hello) Type() Type { return TypeHello }

type ConfigRequest struct{}

func (r ConfigRequest) Type() Type { return TypeConfigRequest }

type ConfigResponse struct {
	Config config.Config
}

func (r ConfigResponse) Type() Type { return TypeConfigResponse }

func Encode(p Interface) []byte {
	buff := new(bytes.Buffer)
	err := gob.NewEncoder(buff).Encode(&p)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}

func Decode(b []byte) Interface {
	var p Interface
	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&p)
	if err != nil {
		panic(err)
	}
	return p
}

func EncodeBase64(p Interface) string {
	b := Encode(p)
	return base64.StdEncoding.EncodeToString(b)
}

func DecodeBase64(s string) Interface {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return Decode(b)
}
