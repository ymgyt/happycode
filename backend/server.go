package backend

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/core/config"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Server) *Server {
	h := cfg.HTTP
	ws := cfg.WebSocket

	s := &Server{
		httpServer: &http.Server{
			Addr:         cfg.Host + ":" + strconv.Itoa(h.Port),
			ReadTimeout:  h.ReadTimeout,
			WriteTimeout: h.WriteTimeout,
			IdleTimeout:  h.IdleTimeout,
		},
		webSocketServer: &http.Server{
			Addr:         cfg.Host + ":" + strconv.Itoa(ws.Port),
			ReadTimeout:  h.ReadTimeout,
			WriteTimeout: h.WriteTimeout,
			IdleTimeout:  h.IdleTimeout,
		},
		errCh: make(chan error, 2), // for each http, websocket server.
	}

	return s
}

type Server struct {
	httpServer      *http.Server
	webSocketServer *http.Server
	errCh           chan error
}

func (s *Server) HTTPHandler(h http.Handler) *Server {
	s.httpServer.Handler = h
	return s
}

func (s *Server) WebSocketHandler(h http.Handler) *Server {
	s.webSocketServer.Handler = h
	return s
}

func (s *Server) Run(ctx context.Context) {
	log.V(0).Info("server running", zap.String("http", s.httpServer.Addr), zap.String("ws", s.webSocketServer.Addr))

	go s.watch(ctx)
	go s.listenAndServeWS()
	s.listenAndServe()
}

func (s *Server) listenAndServe() {
	err := s.httpServer.ListenAndServe()
	if err == http.ErrServerClosed {
		log.V(0).Info("server closed")
	} else {
		log.Error("server stopping", zap.Error(err))
	}
	s.errCh <- err
}

func (s *Server) listenAndServeWS() {
	err := s.webSocketServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Error("websocket server closed", zap.Error(err))
	}
	s.errCh <- err
}

func (s *Server) watch(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-s.errCh:
	}
	bg := context.Background()
	_ = s.httpServer.Shutdown(bg)
	_ = s.webSocketServer.Shutdown(bg)
}
