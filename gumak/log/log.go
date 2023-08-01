package log

import (
	"fmt"
)

type LogLevel int

const (
	LevelNone    LogLevel = 0
	LevelError   LogLevel = 1
	LevelWarning LogLevel = 2
	LevelInfo    LogLevel = 3
	LevelDebug   LogLevel = 4
)

const (
	colorReset  string = "\033[0m"
	colorRed    string = "\033[31m"
	colorGreen  string = "\033[32m"
	colorYellow string = "\033[33m"
	colorBlue   string = "\033[34m"
	colorPurple string = "\033[35m"
	colorCyan   string = "\033[36m"
	colorGray   string = "\033[37m"
	colorWhite  string = "\033[97m"
)

var logLevel LogLevel = LevelInfo

var levelName = []string{
	"",
	"[\033[31merror\033[0m]",
	"[\033[36mwarning\033[0m]",
	"[\033[33minfo\033[0m]",
	"[\033[37mdebug\033[0m]",
	"[\033[35mtrace\033[0m]",
}

func SetLevel(level LogLevel) {
	logLevel = level
}

func Log(level LogLevel, format string, args ...interface{}) {
	if level <= logLevel {
		// TODO: Time & date.
		//fmt.Print(time.Now())
		fmt.Printf("%s ", levelName[level])
		fmt.Printf(format, args...)
		fmt.Print("\n")
	}
}

func Info(fmt string, args ...interface{}) {
	Log(LevelInfo, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	Log(LevelDebug, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	Log(LevelError, fmt, args...)
}

func Warning(fmt string, args ...interface{}) {
	Log(LevelWarning, fmt, args...)
}
