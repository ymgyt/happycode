package frontend

import (
	"fmt"

	"github.com/ymgyt/happycode/frontend/js"
	"github.com/ymgyt/happycode/frontend/service"
)

type App struct {
	d     *js.Element
	world *js.Element

	backendClient *service.BackendClient
	ws            *js.WebSocket
}

func NewApp() *App {
	d := js.NewDocument()
	world, err := d.GetElementByID("world")
	if err != nil {
		panic(err)
	}
	backendURL := js.BackendURL()
	client := service.NewBackendClient(backendURL)

	webSocketPort := client.GetWebSocketPort()
	webSocketEndpoint := fmt.Sprintf("ws://%s:%d", backendURL.Hostname(), webSocketPort)
	ws := js.NewWebSocket(webSocketEndpoint)

	app := &App{
		d:             d,
		world:         world,
		backendClient: client,
		ws:            ws,
	}
	return app
}

func (app *App) Init() {
	s := app.world.Style()

	s.SetWidth("100%")
	s.SetHeight("100%")
	s.SetBackgroundColor("#07280e")
	//s.SetBackgroundColor("#eeeeee")

	app.ws.Init()

	select {}
}
