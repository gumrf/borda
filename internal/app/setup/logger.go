package setup

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
	lock sync.Mutex
}

var ZapLogger = new(Logger)

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

	writerSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(file),
		zapcore.AddSync(os.Stdout),
	)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	ZapLogger.SugaredLogger = sugarLogger
	ZapLogger.lock = sync.Mutex{}

	return nil
}

// Log returns zap logger instance
// Name Log is for convinience
func GetLogger() *zap.SugaredLogger {
	return ZapLogger.SugaredLogger
}
