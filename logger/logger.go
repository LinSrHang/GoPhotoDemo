package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func newLogger() *log.Logger {
	logFile, err := os.Create("./log_" + time.Now().Format("20060102") + ".txt")
	if err != nil {
		fmt.Println(err)
	}

	loger := log.New(logFile, "test_", log.Ldate|log.Ltime|log.Lshortfile)
	return loger
}

func LoggerInit() {
	Logger = newLogger()
}
