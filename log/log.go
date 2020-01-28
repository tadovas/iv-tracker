package log

import (
	"fmt"
	"time"
)

func Error(args ...interface{}) {
	stdPrintln("ERROR", args...)
}

func Info(args ...interface{}) {
	stdPrintln("INFO", args...)
}

func stdPrintln(tag string, args ...interface{}) {
	args = append([]interface{}{time.Now().Format(time.RFC3339), tag}, args...)
	fmt.Println(args...)
}

func IfError(tag string, f func() error) {
	if err := f(); err != nil {
		Error("[", tag, "]", err)
	}
}
