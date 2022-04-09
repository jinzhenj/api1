package log

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Development initialize a development logger
func Development(logLevel int8) logr.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(-zapcore.Level(logLevel))
	return buildLog(cfg)
}

// Production initialize a default logger to be used in production,
// it is used as the default logger.
func Production(logLevel int8) logr.Logger {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(-zapcore.Level(logLevel))
	return buildLog(cfg)
}

func buildLog(cfg zap.Config) logr.Logger {
	zapLog, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return zapr.NewLogger(zapLog)
}
