package vmio

import "fmt"

//VMLoggerConsole use console as a logger platform
type VMLoggerConsole struct {
	level uint16
}

func (log *VMLoggerConsole) Log(level uint16, s string) {
	if level <= log.level {
		fmt.Println(s)
	}
}

func MakeVMLoggerConsole(level uint16) VMLogger {
	return &VMLoggerConsole{level}
}
