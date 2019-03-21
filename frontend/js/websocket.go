package js

import (
	"fmt"
	"syscall/js"

	"github.com/ymgyt/happycode/core/payload"
)

// https://developer.mozilla.org/en-US/docs/Web/API/WebSocket/readyState
type WebSocketState int

const (
	WebSocketConnectiong WebSocketState = 0
	WebSocketOpen        WebSocketState = 1
	WebSocketClosing     WebSocketState = 2
	WebSocketClosed      WebSocketState = 3

	// https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent
	WebSocketCloseEventNormalClosure = 1000
)

func (s WebSocketState) String() string {
	switch s {
	case WebSocketConnectiong:
		return "connecting"
	case WebSocketOpen:
		return "open"
	case WebSocketClosing:
		return "closing"
	case WebSocketClosed:
		return "closed"
	}
	return "unknown"
}

type WebSocket struct {
	ws js.Value
}

func NewWebSocket(endpoint string) *WebSocket {
	ws := js.Global().Get("WebSocket").New(endpoint)
	if isNull(ws) {
		panic("fail: new WebSocket")
	}
	return &WebSocket{ws: ws}
}

func (ws *WebSocket) Init() {
	var cb js.Func
	cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		fmt.Println(fmt.Sprintf("%s: %s", ws.State(), ws.URL()))
		hello := payload.Hello{Message: "hello server !"}
		ws.ws.Call("send", js.TypedArrayOf(payload.Encode(hello)))
		cb.Release()
		return nil
	})
	ws.ws.Call("addEventListener", "open", cb)
}

func (ws *WebSocket) Close() {
	ws.ws.Call("close", WebSocketCloseEventNormalClosure)
}

func (ws *WebSocket) State() WebSocketState {
	return WebSocketState(ws.ws.Get("readyState").Int())
}

func (ws *WebSocket) URL() string {
	return ws.ws.Get("url").String()
}
