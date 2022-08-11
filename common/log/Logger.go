package logger

import (
	"fmt"
	"time"
)

func Info(format string, a ...interface{}) {
	format = getTime() + format
	if len(a) > 0 {
		fmt.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format + "\n")
	}
}

func Debug(format string, a ...interface{}) {
	format = getTime() + format
	if len(a) > 0 {
		fmt.Printf(format+"\n", a...)
	} else {
		fmt.Printf(format + "\n")
	}
}

func Error(format string, a ...interface{}) {
	format = getTime() + format
	fmt.Printf(format+"\n", a...)
}

func Err(format string, a ...interface{}) {
	format = getTime() + format
	fmt.Printf(format+", error: %v\n", a...)
}

func getTime() string {
	return fmt.Sprintf("%v", time.Now())[:19] + "   --->   "
}
