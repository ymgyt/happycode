package handlers

import (
	"encoding/binary"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/payload"
)

type Config struct {
	Config   *config.Config
	Outgoing chan<- payload.Interface
}

func (c *Config) WebSocketPort(w http.ResponseWriter, req *http.Request) {
	err := binary.Write(w, binary.LittleEndian, uint16(c.Config.Server.WebSocket.Port))
	if err != nil {
		log.Error("handle websocket port", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (c *Config) HandleConfig(pl payload.Interface) {
	_, ok := pl.(payload.ConfigRequest)
	if !ok {
		panic(fmt.Errorf("handle config request, got %s, want %s", pl.Type(), payload.TypeConfigRequest))
	}
	c.Outgoing <- payload.ConfigResponse{Config: *c.Config}
}
