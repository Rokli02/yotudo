package logger

import (
	"fmt"
	"time"
)

// Console Logger Date Format
const cldf = "2006-01-02_15:04:05.0000"

type ConsoleLogger uint8

func NewConsoleLogger() Logger {
	return ConsoleLogger(0)
}

var _ Logger = ConsoleLogger(0)

func (c ConsoleLogger) Debug(a ...any) {
	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s  ", time.Now().Format(cldf))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func (c ConsoleLogger) DebugF(format string, a ...any) {
	fmt.Printf("\x1B[38;2;38;161;34m[DEBUG]\t%s  ", time.Now().Format(cldf))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func (c ConsoleLogger) Error(filename string, line int, a ...any) {
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s [%s:%d] ", time.Now().Format(cldf), filename, line)
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func (c ConsoleLogger) ErrorF(filename string, line int, format string, a ...any) {
	fmt.Printf("\x1B[38;2;171;15;18m[ERR]\t%s [%s:%d] ", time.Now().Format(cldf), filename, line)
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func (c ConsoleLogger) Info(a ...any) {
	fmt.Printf("\x1B[38;2;110;195;235m[INFO]\t%s  ", time.Now().Format(cldf))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func (c ConsoleLogger) InfoF(format string, a ...any) {
	fmt.Printf("\x1B[38;2;110;195;235m[INFO]\t%s  ", time.Now().Format(cldf))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func (c ConsoleLogger) Warning(a ...any) {
	fmt.Printf("\x1B[38;2;156;93;39m[WARN]\t%s  ", time.Now().Format(cldf))
	fmt.Println(a...)
	fmt.Printf("\x1B[0m")
}

func (c ConsoleLogger) WarningF(format string, a ...any) {
	fmt.Printf("\x1B[38;2;156;93;39m[WARN]\t%s  ", time.Now().Format(cldf))
	fmt.Printf(format, a...)
	fmt.Printf("\x1B[0m\n")
}

func (c ConsoleLogger) Close() {
	c.Info("Closing logger channels...")
}
