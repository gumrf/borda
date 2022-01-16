package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log contains the shared Logger
	Log *zap.SugaredLogger
)

const (
	_defaultLogLevel = zap.InfoLevel
)

// InitLogger initialize logger with file
func InitLogger(logDir, logFileName string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.Mkdir(logDir, 0777)
		if err != nil {
			return err
		}
	}

	filePath := filepath.Join(logDir, logFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	writerSyncer := zapcore.NewMultiWriteSyncer(zapcore.AddSync(file),
		zapcore.AddSync(os.Stdout))

	encoder := func() zapcore.Encoder {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		return zapcore.NewConsoleEncoder(encoderConfig)
	}()

	core := zapcore.NewCore(encoder, writerSyncer, _defaultLogLevel)

	Log = zap.New(core).Sugar()

	return nil
}
