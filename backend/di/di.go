package di

import (
	"fmt"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/ymgyt/happycode/backend"
	"github.com/ymgyt/happycode/backend/handlers"
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/backend/middlewares"
	"github.com/ymgyt/happycode/core/config"
	"github.com/ymgyt/happycode/core/payload"
)

type AppStack struct {
	Server         *backend.Server
	PayloadManager *backend.PayloadManager
}

func NewAppStack(cfg *config.Config) *AppStack {
	builder := &Builder{cfg: cfg}
	return builder.
		buildPayloadManager().
		buildWebSocket().
		buildRouter().
		buildMWChain().
		buildServer().
		buildAppStack()
}

type Builder struct {
	cfg            *config.Config
	payloadManager *backend.PayloadManager
	websocket      *backend.WebSocket
	router         *backend.Router
	mwChain        *middlewares.Chain
	server         *backend.Server
}

func (b *Builder) buildAppStack() *AppStack {
	return &AppStack{
		Server:         b.server,
		PayloadManager: b.payloadManager,
	}
}

func (b *Builder) buildPayloadManager() *Builder {
	pm := &backend.PayloadManager{
		In:  make(chan payload.Interface, config.PayloadManagerIncommingChanBuffSize),
		Out: make(chan payload.Interface, config.PayloadManagerOutgoingChanBuffSize),
	}
	hello := &handlers.Hello{Outgoing: pm.Out}
	pm.Register(payload.TypeHello, hello.HandleHello)

	b.payloadManager = pm
	return b
}

func (b *Builder) buildWebSocket() *Builder {
	origin := fmt.Sprintf("http://%s:%d", b.cfg.Server.Host, b.cfg.Server.HTTP.Port)
	wsCfg := b.cfg.Server.WebSocket
	ws := &backend.WebSocket{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  wsCfg.ReadBufferSize,
			WriteBufferSize: wsCfg.WriteBufferSize,
			CheckOrigin:     backend.WebSocketCheckOrigin([]string{origin}),
		},
	}
	ws.IncommingPayload = b.payloadManager.In
	ws.OutgoingPayload = b.payloadManager.Out
	b.websocket = ws
	return b
}

func (b *Builder) buildRouter() *Builder {
	r := backend.NewRouter()
	staticDir := b.cfg.Server.StaticDir()
	staticHdler := handlers.NewStatic(staticDir, "/"+filepath.Base(staticDir))
	configHdler := handlers.NewConfig(b.cfg)

	r.GET("/", staticHdler.ServeHTTP)
	r.GET("/favicon.ico", staticHdler.ServeHTTP)
	r.GET("/static/*", staticHdler.ServeHTTP)
	r.GET("/config/server.websocket.port", configHdler.WebSocketPort)
	b.router = r
	return b
}

func (b *Builder) buildMWChain() *Builder {
	chain := middlewares.NewChain(b.router, []middlewares.Interface{
		middlewares.MustLogging(&middlewares.LoggingConfig{Logger: log.Clone(), Console: true}),
	})
	b.mwChain = chain
	return b
}

func (b *Builder) buildServer() *Builder {
	server := backend.NewServer(b.cfg.Server).
		HTTPHandler(b.mwChain).
		WebSocketHandler(b.websocket)
	b.server = server
	return b
}
