package logger

import (
	"fmt"
	"runtime"
	"time"
)

const isDebugMode = true
const dateFormat = "2006-01-02_15:04:05.0000"

type Logger interface {
	Info(a ...any)
	InfoF(format string, a ...any)
	Debug(a ...any)
	DebugF(format string, a ...any)
	WarningF(format string, a ...any)
	Warning(a ...any)
	ErrorF(format string, a ...any)
	Error(a ...any)
}

func InitializeLogger() ([]Logger, error) {

	return nil, nil
}

func Info(a ...any) {
	fmt.Printf("\x1B[38;2;110;195;235m[INFO]\t%s  ", time.Now().Format(dateFormat))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func InfoF(format string, a ...any) {
	fmt.Printf("\x1B[38;2;110;195;235m[INFO]\t%s  ", time.Now().Format(dateFormat))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func Debug(a ...any) {
	if !isDebugMode {
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s [%s:%d]  ", time.Now().Format(dateFormat), filename, line)
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func DebugF(format string, a ...any) {
	if !isDebugMode {
		return
	}

	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s [%s:%d]  ", time.Now().Format(dateFormat), filename, line)
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func WarningF(format string, a ...any) {
	fmt.Printf("\x1B[38;2;156;93;39m[WARN]\t%s  ", time.Now().Format(dateFormat))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func Warning(a ...any) {
	fmt.Printf("\x1B[38;2;156;93;39m[WARN]\t%s  ", time.Now().Format(dateFormat))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func ErrorF(format string, a ...any) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s [%s:%d] ", time.Now().Format(dateFormat), filename, line)
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func Error(a ...any) {
	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s [%s:%d] ", time.Now().Format(dateFormat), filename, line)
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}
