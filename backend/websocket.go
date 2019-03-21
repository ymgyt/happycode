package backend

import (
	"bytes"
	"encoding/gob"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/payload"
	_ "github.com/ymgyt/happycode/core/payload"
	"go.uber.org/zap"
)

type WebSocket struct {
	upgrader websocket.Upgrader
}

func NewWebSocket(cfg *config.WebSocket) *WebSocket {
	return &WebSocket{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.ReadBufferSize,
			WriteBufferSize: cfg.WriteBufferSize,
			CheckOrigin:     WebSocketCheckOrigin,
		},
	}
}

func WebSocketCheckOrigin(req *http.Request) bool {
	log.V(10).Debug("websocket check origin", zap.Any("req", req))
	return true
}

func (ws *WebSocket) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Error("websocket upgrade connection", zap.Error(err))
		return
	}
	wsConn := &WebSocketConn{Conn: conn}
	go wsConn.read()
	// go ws.write()
}

type WebSocketConn struct {
	*websocket.Conn
}

func (conn *WebSocketConn) read() {
	defer conn.Close()

	// TODO Set Read/Write Timeout
	conn.SetPongHandler(func(s string) error {
		// TODO update readDeadline
		log.V(10).Debug("websocket", zap.String("pong", s))
		return nil
	})

	for {
		/*
					// The message types are defined in RFC 6455, section 11.8.
			const (
				// TextMessage denotes a text data message. The text message payload is
				// interpreted as UTF-8 encoded text data.
				TextMessage = 1

				// BinaryMessage denotes a binary data message.
				BinaryMessage = 2

				// CloseMessage denotes a close control message. The optional message
				// payload contains a numeric code and text. Use the FormatCloseMessage
				// function to format a close message payload.
				CloseMessage = 8

				// PingMessage denotes a ping control message. The optional message payload
				// is UTF-8 encoded text.
				PingMessage = 9

				// PongMessage denotes a pong control message. The optional message payload
				// is UTF-8 encoded text.
				PongMessage = 10
			)
		*/

		messageType, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("websocket", zap.String("action", "readMessage"), zap.Error(err))
				break
			}
		}
		log.Debug("websocket", zap.Int("message_type", messageType), zap.ByteString("data", data))

		var payload payload.Interface
		err = gob.NewDecoder(bytes.NewReader(data)).Decode(&payload)
		if err != nil {
			log.Error("websocket", zap.String("action", "gob decode"), zap.Error(err))
			break
		}

		spew.Dump(payload)
	}
}
