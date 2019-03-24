package frontend

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/davecgh/go-spew/spew"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/payload"
	"github.com/ymgyt/happycode/frontend/js"
	"github.com/ymgyt/happycode/frontend/log"
	"github.com/ymgyt/happycode/frontend/service"
	"github.com/ymgyt/happycode/frontend/ui"
)

type App struct {
	// Document *js.Element
	// World    *js.Element
	UI *ui.UI

	BackendClient *service.BackendClient
	WebSocket     *js.WebSocket
	Config        *config.Config
}

func (app *App) Run() {
	log.V(10).Debug("initialize app")
	app.Init()
	log.V(10).Debug("running app")
	select {}
}

func (app *App) Init() {
	app.WebSocket.Init()
	go app.HandlePayload()
	go app.HandleKeyboardEvent()

	app.LoadConfig()
}

func (app *App) LoadConfig() {
	app.WebSocket.Send(payload.ConfigRequest{})
}

func (app *App) HandleKeyboardEvent() {
	for ke := range app.UI.KeyboardEvents {
		spew.Dump(ke)
	}
}

func (app *App) HandlePayload() {
	for pl := range app.WebSocket.IncommingPayload {
		typ := pl.Type()
		log.V(0).Debug("receive payload", zap.String("type", typ.String()))

		switch typ {
		case payload.TypeHello:
			app.HandleHello(pl)
		case payload.TypeConfigResponse:
			app.HandleConfigResponse(pl)
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

func (app *App) HandleConfigResponse(pl payload.Interface) {
	cfgResp, ok := pl.(payload.ConfigResponse)
	if !ok {
		panic(fmt.Errorf("invalid payload type, got %s, want %s", pl.Type(), payload.TypeConfigResponse))
	}
	// Note: should i care concurrency ?
	app.Config = &cfgResp.Config
	log.V(0).Info("load config")
	app.UI.Init(app.Config)
}

