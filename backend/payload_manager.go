package backend

import (
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/payload"
	"go.uber.org/zap"
)

type PayloadManager struct {
	In  chan payload.Interface
	Out chan payload.Interface

	entries map[payload.Type][]PayloadHandler
}

func (pm *PayloadManager) Run() {
	go pm.read()
}

func (pm *PayloadManager) Register(t payload.Type, h PayloadHandler) {
	if pm.entries == nil {
		pm.entries = make(map[payload.Type][]PayloadHandler)
	}
	pm.entries[t] = append(pm.entries[t], h)
}

func (pm *PayloadManager) read() {
	for {
		pl := <-pm.In
		for _, handler := range pm.lookupHandlers(pl.Type()) {
			handler(pl)
		}
	}
}

type PayloadHandler func(pl payload.Interface)

func (pm *PayloadManager) lookupHandlers(typ payload.Type) []PayloadHandler {
	hds, found := pm.entries[typ]
	if !found {
		return []PayloadHandler{pm.notFound}
	}
	return hds
}

func (pm *PayloadManager) notFound(pl payload.Interface) {
	log.V(0).Error("handle unregister payload", zap.String("payload_type", pl.Type().String()))
}
