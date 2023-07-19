package tlog

import (
	"fmt"
	"log"
	"os"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var level = LevelInfo

func SetLevel(l int) {
	level = l
}

func GetLevel() int {
	return level
}

var TLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

func SetLog(l *log.Logger) {
	TLog = l
}

func Debug(fmtS string, v ...interface{}) {
	if level <= LevelDebug {
		TLog.Printf(fmt.Sprintf("[DEBUG] %s\n", fmtS), v...)
	}
}

func Info(fmtS string, v ...interface{}) {
	if level <= LevelInfo {
		TLog.Printf(fmt.Sprintf("[INFO] %s\n", fmtS), v...)
	}
}

func Warn(fmtS string, v ...interface{}) {
	if level <= LevelWarn {
		TLog.Printf(fmt.Sprintf("[WARN] %s\n", fmtS), v...)
	}
}

func Error(fmtS string, v ...interface{}) {
	if level <= LevelError {
		TLog.Printf(fmt.Sprintf("[ERROR] %s\n", fmtS), v...)
	}
}
