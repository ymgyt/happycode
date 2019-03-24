package di

import (
	"fmt"

	"github.com/ymgyt/happycode/core/payload"
	"github.com/ymgyt/happycode/frontend"
	"github.com/ymgyt/happycode/frontend/js"
	"github.com/ymgyt/happycode/frontend/service"
	"github.com/ymgyt/happycode/frontend/ui"
)

func NewApp() *frontend.App {
	backendURL := js.BackendURL()
	client := service.NewBackendClient(backendURL)

	webSocketPort := client.GetWebSocketPort()
	webSocketEndpoint := fmt.Sprintf("ws://%s:%d", backendURL.Hostname(), webSocketPort)
	ws := js.NewWebSocket(webSocketEndpoint)

	ws.IncommingPayload = make(chan payload.Interface, 100)
	ws.OutgoingPayload = make(chan payload.Interface, 100)

	app := &frontend.App{
		UI:            NewUI(),
		BackendClient: client,
		WebSocket:     ws,
	}
	return app
}

func NewUI() *ui.UI {
	d := js.NewDocument()
	world, err := d.GetElementByID("world")
	if err != nil {
		panic(err)
	}
	return &ui.UI{
		Document:       d,
		World:          world,
		KeyboardEvents: make(chan js.KeyboardEvent, 100),
		ActionEvents:   make(chan js.KeyboardEvent, 100),
	}
}
