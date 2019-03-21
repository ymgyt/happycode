package di

import (
	"path/filepath"

	"github.com/ymgyt/happycode/backend"
	"github.com/ymgyt/happycode/backend/handlers"
	"github.com/ymgyt/happycode/backend/log"
	"github.com/ymgyt/happycode/backend/middlewares"
	"github.com/ymgyt/happycode/core/config"
)

type AppStack struct {
	Server *backend.Server
}

func NewAppStack(cfg *config.Config) *AppStack {
	return &AppStack{
		Server: NewServer(cfg),
	}
}

func NewServer(cfg *config.Config) *backend.Server {

	router := NewRouter(cfg)
	chain := middlewares.NewChain(router, []middlewares.Interface{
		middlewares.MustLogging(&middlewares.LoggingConfig{Logger: log.Clone(), Console: true}),
	})
	ws := backend.NewWebSocket(cfg.Server.WebSocket)

	server := backend.NewServer(cfg.Server).
		HTTPHandler(chain).
		WebSocketHandler(ws)

	return server
}

func NewRouter(cfg *config.Config) *backend.Router {
	r := backend.NewRouter()
	hg := newHandlerGroup(cfg)

	r.GET("/", hg.static.ServeHTTP)
	r.GET("/favicon.ico", hg.static.ServeHTTP)
	r.GET("/static/*", hg.static.ServeHTTP)
	r.GET("/config/server.websocket.port", hg.config.WebSocketPort)

	return r
}

type handlerGroup struct {
	static *handlers.Static
	config *handlers.Config
}

func newHandlerGroup(cfg *config.Config) *handlerGroup {
	staticDir := cfg.Server.StaticDir()
	return &handlerGroup{
		static: handlers.NewStatic(staticDir, "/"+filepath.Base(staticDir)),
		config: handlers.NewConfig(cfg),
	}
}
