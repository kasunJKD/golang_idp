package logger

import (
	"log"
	"os"
)

type CustomLogger struct {
    Debuglog *log.Logger
    Infolog   *log.Logger
    Errorlog   *log.Logger
}

func NewCustomLogger() *CustomLogger {
    return &CustomLogger{
        Infolog: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
        Debuglog: log.New(os.Stdout, "Debug: ", log.Ldate|log.Ltime|log.Lshortfile),
        Errorlog: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}