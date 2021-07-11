package log

import (
	"fmt"
	"time"
)

// ANSI colors: https://gist.github.com/JBlond/2fea43a3049b38287e5e9cefc87b2124
var (
	None   = "0"
	Red    = "31"
	Green  = "32"
	Yellow = "33"
	Blue   = "34"
)

// Debug logs with tag debug
func Debug(msg string) {
	print(msg, "DEBUG", None)
}

// Info logs with tag info and in blue
func Info(msg string) {
	print(msg, "INFO", Blue)
}

// Warn logs with tag warn and in yellow
func Warn(msg string) {
	print(msg, "WARNING", Yellow)
}

// Error logs with tag error and in red
func Error(msg string) {
	print(msg, "ERROR", Red)
}

func print(msg string, severity string, color string) {
	currentTime := time.Now()
	fmt.Printf("\033[0;%sm%s - [%s] --> %s\033[0m\n", color, currentTime.Format("2006-01-02 15:04:05"), severity, msg)
}
