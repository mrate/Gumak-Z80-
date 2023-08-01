//go:build !trace

package log

func Trace(level int, fmt string, args ...interface{}) {}
func TraceEnable(level int)                            {}
func TraceFlush()                                      {}
