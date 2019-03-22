package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// config file
	happyCodeRootDirName = ".happycode"
	configFileName       = "config.yaml"
	staticDirName        = "static"

	// server
	defaultServerRootDir    = "${GOPATH}/src/github.com/ymgyt/happycode"
	defaultServerHost       = "localhost"
	defaultHTTPPort         = 8888
	defaultHTTPReadTimeout  = 3 * time.Second
	defaultHTTPWriteTimeout = 3 * time.Second
	defaultHTTPIdleTimeout  = 3 * time.Second

	// websocket
	defaultWebSocketPort            = 50505
	defaultWebSocketReadBufferSize  = 4096
	defaultWebSocketWriteBufferSize = 4096

	// theme
	defaultThemeBackgroundColor = "#07280e"
)

const (
	// payload manager
	PayloadManagerIncommingChanBuffSize = 100
	PayloadManagerOutgoingChanBuffSize  = 100
)

var (
	happyCodeRootDirPermission = os.ModeDir | OS_USER_RWX
	configFilePermission       = OS_USER_RW
)

func Default() *Config {
	cfg := &Config{}
	cfg.SetDefault()
	return cfg
}

type Config struct {
	Meta   *Meta `yaml:"-"`
	Server *Server
	Theme  *Theme
}

type Meta struct {
	FilePath string
}

type Server struct {
	Host      string
	RootDir   string
	HTTP      *HTTP
	WebSocket *WebSocket
}

type HTTP struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type WebSocket struct {
	Port            int
	ReadBufferSize  int
	WriteBufferSize int
}

type Theme struct {
	BackgroundColor string
}

func NewFromYAML(b []byte, path string) (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(b, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.Meta = &Meta{FilePath: path}
	cfg.SetDefault()
	return &cfg, nil
}

func (c *Config) SetDefault() {
	if c.Meta == nil {
		c.Meta = &Meta{}
	}
	c.Meta.SetDefault()

	if c.Server == nil {
		c.Server = &Server{}
	}
	c.Server.SetDefault()

	if c.Theme == nil {
		c.Theme = &Theme{}
	}
	c.Theme.SetDefault()
}

func (c *Config) FilePermission() os.FileMode { return os.FileMode(configFilePermission) }
func (c *Config) DirPermission() os.FileMode  { return os.FileMode(happyCodeRootDirPermission) }

func (m *Meta) SetDefault() {}

func (s *Server) StaticDir() string {
	return filepath.Join(s.RootDir, staticDirName)
}
func (s *Server) SetDefault() {
	if s.Host == "" {
		s.Host = defaultServerHost
	}
	if s.RootDir == "" {
		s.RootDir = os.ExpandEnv(defaultServerRootDir)
	}
	if s.HTTP == nil {
		s.HTTP = &HTTP{}
	}
	if s.WebSocket == nil {
		s.WebSocket = &WebSocket{}
	}
	s.HTTP.SetDefault()
	s.WebSocket.SetDefault()
}

func (h *HTTP) SetDefault() {
	if h.Port == 0 {
		h.Port = defaultHTTPPort
	}
	if h.ReadTimeout == 0 {
		h.ReadTimeout = defaultHTTPReadTimeout
	}
	if h.WriteTimeout == 0 {
		h.WriteTimeout = defaultHTTPWriteTimeout
	}
	if h.IdleTimeout == 0 {
		h.IdleTimeout = defaultHTTPIdleTimeout
	}
}

func (w *WebSocket) SetDefault() {
	if w.Port == 0 {
		w.Port = defaultWebSocketPort
	}
	if w.ReadBufferSize == 0 {
		w.ReadBufferSize = defaultWebSocketReadBufferSize
	}
	if w.WriteBufferSize == 0 {
		w.WriteBufferSize = defaultWebSocketWriteBufferSize
	}
}

func (t *Theme) SetDefault() {
	if t.BackgroundColor == "" {
		t.BackgroundColor = defaultThemeBackgroundColor
	}
}

func DefaultDir() string {
	home := HomeDir()
	return filepath.Join(home, happyCodeRootDirName)
}

func DefaultConfigPath() string {
	return filepath.Join(DefaultDir(), configFileName)
}

func HomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic("can not determine user home dir " + err.Error())
	}
	return home
}

// stole from https://stackoverflow.com/questions/28969455/golang-properly-instantiate-os-filemode
const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)
