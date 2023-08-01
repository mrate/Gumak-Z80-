//go:build trace

package log

import "fmt"

var tracingLevel = 0
var dirty = false

func Trace(level int, format string, args ...interface{}) {
	if tracingLevel >= level {
		fmt.Printf("\033[35m")
		fmt.Printf(format, args...)
		fmt.Printf("\033[0m")

		dirty = true
	}
}

func TraceEnable(level int) {
	tracingLevel = level
}

func TraceFlush() {
	if dirty {
		fmt.Print("\n")
		dirty = false
	}
}
