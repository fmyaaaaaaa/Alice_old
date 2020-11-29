package logger

import (
	"log"
	"os"
	"sync"
)

type AliceLogManager struct {
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

var logManager *AliceLogManager
var once sync.Once

func LogManager() AliceLogManager {
	once.Do(func() {
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		logManager = &AliceLogManager{}
		var manager = AliceLogManager{
			WarningLogger: log.New(file, "[Warning] : ", log.Ldate|log.Ltime|log.Lshortfile),
			InfoLogger:    log.New(file, "[Info] : ", log.Ldate|log.Ltime|log.Lshortfile),
			ErrorLogger:   log.New(file, "[Error] : ", log.Ldate|log.Ltime|log.Lshortfile),
		}
		*logManager = manager
	})
	return *logManager
}

func (l AliceLogManager) Warning(message ...interface{}) {
	l.WarningLogger.Println(message)
}

func (l AliceLogManager) Info(message ...interface{}) {
	l.InfoLogger.Println(message)
}

func (l AliceLogManager) Error(message ...interface{}) {
	l.ErrorLogger.Println(message)
}
