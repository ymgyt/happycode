package handlers

import (
	"encoding/binary"
	"net/http"

	"go.uber.org/zap"

	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
)

type Config struct {
	cfg *config.Config
}

func NewConfig(cfg *config.Config) *Config {
	return &Config{cfg: cfg}
}

func (c *Config) WebSocketPort(w http.ResponseWriter, req *http.Request) {
	err := binary.Write(w, binary.LittleEndian, uint16(c.cfg.Server.WebSocket.Port))
	if err != nil {
		log.Error("handle websocket port", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
