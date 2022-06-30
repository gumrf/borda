package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogLevel = zap.InfoLevel
)

// Logger is a wrapper around zap.ShugaredLogger
type Logger struct {
	*zap.SugaredLogger
}

type Config struct {
	Directory string
	Outputs   []string
}

func New() *Logger {
	config := Config{
		Outputs: []string{"stderr", "./logs/borda.log"},
	}

	zapConfig := zap.Config{
		Level:    zap.NewAtomicLevelAt(defaultLogLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: LevelEncoder,

			TimeKey:    "time",
			EncodeTime: TimeEncoder,
		},
		OutputPaths: config.Outputs,
	}

	baseLogger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	baseLogger.Info("Logger Initialized")

	return &Logger{
		SugaredLogger: baseLogger.Sugar(),
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func LevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}
