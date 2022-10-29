package logger

import (
	"log"
	"os"
)

var (
	DebugLogger   *log.Logger
	ErrorLogger   *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
)

func init() {
	ErrorLogger = log.New(os.Stderr, "FLUFIK ERROR: 🔥", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	DebugLogger = log.New(os.Stdout, "FLUFIK DEBUG: ❕", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	WarningLogger = log.New(os.Stdout, "FLUFIK WARNING: ⚠️", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
	InfoLogger = log.New(os.Stdout, "FLUFIK INFO: ✅", log.Ldate|log.Ltime|log.Lmsgprefix|log.LUTC|log.Lmicroseconds)
}

func RaiseErr(message string, err ...interface{}) {
	ErrorLogger.Printf(message, err...)
	os.Exit(1)
}

func RaiseWarn(message string, warn ...interface{}) {
	WarningLogger.Printf(message, warn...)
}

func InfoLog(message string, log ...interface{}) {
	InfoLogger.Printf(message, log...)
}

func CheckErr(message string, err error) {
	if err != nil {
		if len(message) == 0 {
			ErrorLogger.Printf("%s\n", err)
			os.Exit(1)
		}
		ErrorLogger.Printf("%s - %s\n", message, err)
		os.Exit(1)
	}
}

func DebugLog(message string, debug ...interface{}) {
	if IsDebugEnabled() {
		DebugLogger.Printf(message, debug...)
	}
}

func IsDebugEnabled() bool {
	if os.Getenv("PERF_DEBUG") == "1" {
		return true
	}
	return false
}
