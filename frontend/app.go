package frontend

import (
	"github.com/ymgyt/happycode/core/payload"
	"github.com/ymgyt/happycode/frontend/js"
	"github.com/ymgyt/happycode/frontend/log"
	"github.com/ymgyt/happycode/frontend/service"
	"go.uber.org/zap"
)

type App struct {
	Document *js.Element
	World    *js.Element

	BackendClient *service.BackendClient
	WebSocket     *js.WebSocket
}

func (app *App) Init() {
	s := app.World.Style()
	s.SetWidth("100%")
	s.SetHeight("100%")
	s.SetBackgroundColor("#07280e")

	app.WebSocket.Init()
	go app.HandlePayload()
}

func (app *App) Run() {
	log.V(10).Debug("initialize app")
	app.Init()
	log.V(10).Debug("running app")
	select {}
}

func (app *App) HandlePayload() {
	for pl := range app.WebSocket.IncommingPayload {
		typ := pl.Type()
		log.V(0).Debug("receive payload", zap.String("type", typ.String()))

		switch typ {
		case payload.TypeHello:
			app.HandleHello(pl)
		default:
			panic("unexpected payload type " + typ.String())
		}
	}
}

func (app *App) HandleHello(pl payload.Interface) {
	hello, ok := pl.(payload.Hello)
	if !ok {
		panic("payload is not hello type")
	}
	log.V(0).Info("handle hello payload", zap.String("message", hello.Message))
}

