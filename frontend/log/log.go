package log

import (
	"os"

	"github.com/ymgyt/happycode/backend/log"
	"go.uber.org/zap"
)

var (
	Verbosity = 0

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

// TODO: move log to core, or create frontend specific logger.
func newDefaultLogger() *zap.Logger {
	cfg := &log.Config{
		Out:    os.Stdout,
		Level:  "debug",
		Encode: log.EncodeConsole,
		Color:  false,
	}
	return log.Must(cfg)
}
