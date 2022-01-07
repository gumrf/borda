package setup

import (
	"fmt"
	"path/filepath"
	"log"
	"os"
	"sync"
)

type LOGGER struct {
    Log *log.Logger
}
var lock = &sync.Mutex{}
var loggers *LOGGER

	
func GetLoggerInstance() *LOGGER {
	lock.Lock()
	defer lock.Unlock()
	if loggers == nil {
		// fmt.Println("Creating LOGGER instance now.")
		loggers = &LOGGER{}
	} // else {
	// 	// fmt.Println("LOGGER instance already created.")
    // }
    return loggers
}

func GetReadyLogger() (*log.Logger, error) {
	lock.Lock()
	defer lock.Unlock()
	if loggers == nil || loggers.Log == nil {
		return nil, fmt.Errorf("logger is not init")
	}
    return loggers.Log, nil
}

func InitLogger(logDir string, logFileName string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err:= os.Mkdir(logDir, 0777)
		if err != nil {
			return err
		}
	}
	filePath := filepath.Join(logDir, logFileName)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err:= os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        return err
    }
	logger:= GetLoggerInstance()
	logger.Log = log.New(file, "INFO\t", log.Ldate|log.Ltime)
	return nil
}