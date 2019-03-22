package handlers

import (
	"fmt"

	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/payload"
	"go.uber.org/zap"
)

type Hello struct {
	Outgoing chan<- payload.Interface
}

func (h *Hello) HandleHello(pl payload.Interface) {
	// not pointer type
	hello, ok := pl.(payload.Hello)
	if !ok {
		log.V(0).Fatal("handle hello payload", zap.Error(fmt.Errorf("invalid payload type. got %s, want %s", pl.Type(), payload.TypeHello)))
	}

	got := hello.Message
	log.V(10).Info("got hello payload", zap.String("msg", got))
	msg := fmt.Sprintf("hello client, I got [%s]", got)
	h.Outgoing <- &payload.Hello{Message: msg}
}
