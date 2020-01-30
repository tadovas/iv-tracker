package log

import (
	"fmt"
	"time"
)

func Error(args ...interface{}) {
	stdPrintln("ERROR", args...)
}

func Errorf(format string, args ...interface{}) {
	stdPrintf("ERROR", format, args...)
}

func Info(args ...interface{}) {
	stdPrintln("INFO", args...)
}

func Infof(format string, args ...interface{}) {
	stdPrintf("INFO", format, args...)
}

func stdPrintf(tag, format string, args ...interface{}) {
	format = "%v %v " + format
	args = append([]interface{}{time.Now().Format(time.RFC3339), tag}, args...)
	fmt.Printf(format, args...)
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
