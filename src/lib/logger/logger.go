package logger

import (
	"fmt"
	"time"
)

const isDebugMode = true
const dateFormat = "2006-01-02_15:04:05.0000"

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

	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s  ", time.Now().Format(dateFormat))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func DebugF(format string, a ...any) {
	if !isDebugMode {
		return
	}

	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s  ", time.Now().Format(dateFormat))
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
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s  ", time.Now().Format(dateFormat))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func Error(a ...any) {
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s  ", time.Now().Format(dateFormat))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}
