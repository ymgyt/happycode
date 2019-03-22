package js

import (
	"syscall/js"

	"github.com/ymgyt/happycode/core/payload"
	"github.com/ymgyt/happycode/frontend/log"
	"go.uber.org/zap"
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
	ws               js.Value
	IncommingPayload chan payload.Interface
	OutgoingPayload  chan payload.Interface
}

func NewWebSocket(endpoint string) *WebSocket {
	ws := js.Global().Get("WebSocket").New(endpoint)
	if isNull(ws) {
		panic("fail: new WebSocket")
	}
	return &WebSocket{ws: ws}
}

func (ws *WebSocket) Init() {
	ws.initOnMessage() // listen on message before send hello.
	ws.initOnOpen()
	ws.initOnError()
	ws.initOnClose()
	go ws.write()
}

func (ws *WebSocket) initOnOpen() {
	var cb js.Func
	cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		log.V(0).Debug("connect websocket", zap.String("state", ws.State().String()), zap.String("endpoint", ws.URL()))
		hello := payload.Hello{Message: "hello server !"}
		ws.Send(hello)
		cb.Release()
		return nil
	})
	ws.ws.Call("addEventListener", "open", cb)
}

func (ws *WebSocket) initOnMessage() {
	ws.ws.Call("addEventListener", "message", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		messageEvent := args[0]
		data := messageEvent.Get("data")
		pl := payload.DecodeBase64(data.String())
		ws.IncommingPayload <- pl
		return nil
	}))
}

func (ws *WebSocket) initOnError() {
	ws.ws.Call("addEventListener", "error", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		err := args[0]
		log.V(0).Error("receive websocket error", zap.String("event", err.String()))
		return nil
	}))
}

func (ws *WebSocket) initOnClose() {
	ws.ws.Call("addEventListener", "close", js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		event := args[0]
		log.V(0).Warn("receive websocket close event", zap.String("event", event.String()))
		return nil
	}))
}

func (ws *WebSocket) Send(pl payload.Interface) {
	ws.OutgoingPayload <- pl
}

func (ws *WebSocket) write() {
	for pl := range ws.OutgoingPayload {
		b := js.TypedArrayOf(payload.Encode(pl))
		ws.ws.Call("send", b)
		b.Release()
	}
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
