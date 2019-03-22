package backend

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/payload"
	_ "github.com/ymgyt/happycode/core/payload"
	"go.uber.org/zap"
)

type WebSocketMessageType int

// The message types are defined in RFC 6455, section 11.8.
const (
	WebSocketTextMessage   WebSocketMessageType = 1
	WebSocketBinaryMessage WebSocketMessageType = 2
	WebSocketCloseMessage  WebSocketMessageType = 8
	WebSocketPingMessage   WebSocketMessageType = 9
	WebSocketPongMessage   WebSocketMessageType = 10
)

func (t WebSocketMessageType) String() string {
	switch t {
	case WebSocketTextMessage:
		return "text"
	case WebSocketBinaryMessage:
		return "binary"
	case WebSocketCloseMessage:
		return "close"
	case WebSocketPingMessage:
		return "ping"
	case WebSocketPongMessage:
		return "pong"
	}
	return "undefined"
}

func (t WebSocketMessageType) Int() int { return int(t) }

type WebSocket struct {
	Upgrader         websocket.Upgrader
	IncommingPayload chan<- payload.Interface
	OutgoingPayload  <-chan payload.Interface
}

func WebSocketCheckOrigin(origins []string) func(*http.Request) bool {
	return func(req *http.Request) bool {
		var ok bool
		var origin = req.Header.Get("Origin")
		for _, allowed := range origins {
			if origin == allowed {
				ok = true
				break
			}
		}
		if ok {
			log.V(10).Debug("websocket check origin", zap.Bool("ok", ok), zap.String("origin", origin))
		} else {
			log.Warn("websocket check origin", zap.Bool("ok", ok), zap.String("origin", origin))
		}
		return ok
	}
}

func (ws *WebSocket) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("websocket upgrade connection", zap.Error(err))
		return
	}
	wsConn := &WebSocketConn{
		Conn:      conn,
		incomming: ws.IncommingPayload,
		outgoing:  ws.OutgoingPayload,
		closed:    make(chan struct{}, 1),
	}
	go wsConn.read()
	go wsConn.write()
}

type WebSocketConn struct {
	*websocket.Conn
	closed    chan struct{}
	incomming chan<- payload.Interface
	outgoing  <-chan payload.Interface
}

func (conn *WebSocketConn) readMessage() (WebSocketMessageType, []byte, error) {
	mt, data, err := conn.Conn.ReadMessage()
	return WebSocketMessageType(mt), data, err
}

func (conn *WebSocketConn) read() {
	defer conn.Close()

	// currently, Javascript API does not provide ping API.
	// https://stackoverflow.com/questions/10585355/sending-websocket-ping-pong-frame-from-browser
	conn.SetPongHandler(func(s string) error {
		log.V(10).Debug("websocket", zap.String("pong", s))
		return nil
	})

	for {
		messageType, data, err := conn.readMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.V(0).Error("websocket", zap.String("action", "readMessage"), zap.Error(err))
				break
			}

			log.V(0).Info("websocket closed", zap.String("reason", err.Error()))
			break
		}
		log.V(10).Debug("read websocket message", zap.String("message_type", messageType.String()))

		pl := payload.Decode(data)
		conn.incomming <- pl
	}
}

func (conn *WebSocketConn) Close() {
	conn.closed <- struct{}{}
	conn.Conn.Close()
}

func (conn *WebSocketConn) write() {
	for {
		select {
		case <-conn.closed:
			return
		case pl := <-conn.outgoing:
			log.V(10).Debug("write payload to websocket", zap.String("type", pl.Type().String()))
			err := conn.Conn.WriteMessage(WebSocketTextMessage.Int(), []byte(payload.EncodeBase64(pl)))
			if err != nil {
				log.V(0).Error("write payload to websocket", zap.String("type", pl.Type().String()), zap.Error(err))
			}
		}
	}
}
