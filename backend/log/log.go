package log

import (
	"io"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EncodeConsole = "console"
	EncodeJSON    = "json"

	defaultVerbosity = 0
)

var (
	Verbosity = defaultVerbosity

	defaultLogger = newDefaultLogger()
	nopLogger     = zap.NewNop()
)

func Debug(k string, fields ...zap.Field) {
	defaultLogger.Debug(k, fields...)
}

func Info(k string, fields ...zap.Field) {
	defaultLogger.Info(k, fields...)
}

func Warn(k string, fields ...zap.Field) {
	defaultLogger.Warn(k, fields...)
}

func Error(k string, fields ...zap.Field) {
	defaultLogger.Error(k, fields...)
}

func Fatal(k string, fields ...zap.Field) {
	defaultLogger.Fatal(k, fields...)
}

func V(v int) *zap.Logger {
	if v < Verbosity {
		return nopLogger
	}
	return defaultLogger
}

func Clone() *zap.Logger {
	return defaultLogger.With()
}

func newDefaultLogger() *zap.Logger {
	cfg := &Config{
		Out:    os.Stdout,
		Level:  "debug",
		Encode: EncodeConsole,
		Color:  true,
	}
	return Must(cfg)
}

// Config -
type Config struct {
	Out       io.Writer
	Level     string
	Encode    string
	Color     bool
	Timestamp bool
}

func (c *Config) timeKey() string {
	if c.Timestamp {
		return "t"
	}
	return ""
}

// Must -
func Must(cfg *Config) *zap.Logger {
	z, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return z
}

// New -
func New(cfg *Config) (*zap.Logger, error) {

	encCfg := zapcore.EncoderConfig{
		TimeKey:        cfg.timeKey(),
		LevelKey:       "l",
		CallerKey:      "c",
		MessageKey:     "m",
		StacktraceKey:  "s",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var level zap.AtomicLevel
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	if cfg.Color {
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Encode == EncodeConsole {
		encoder = zapcore.NewConsoleEncoder(encCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(cfg.Out), level)
	z := zap.New(core, zap.AddCaller())

	return z, nil
}
