package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// File Logger Date Format
	fldf = "15:04:05.0000"
	// File Logger Filename Date Format
	flfdf = "2006-01-02"
	// File Logger's base directory
	baseDir = "./data/logs"
)

type FileLogger struct {
	logFile       *os.File
	writer        *bufio.Writer
	currentTime   int
	stopAutoflush chan bool
}

func NewFileLogger() *FileLogger {
	f := &FileLogger{
		stopAutoflush: make(chan bool),
	}

	filename := filepath.Join(baseDir, fmt.Sprintf("logs_%s.txt", time.Now().Format(flfdf)))
	f.currentTime = f.createCurrentTime()

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, os.ModeType)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			panic("Couldn't create log-file: " + err.Error())
		}
	}

	f.logFile = file
	f.writer = bufio.NewWriter(file)

	go func() {
		ticker := time.NewTicker(time.Second * 10)

		for {
			select {
			case <-ticker.C:
				if f.writer != nil {
					if f.writer.Buffered() != 0 {
						f.writer.Flush()
					}
				}
			case <-f.stopAutoflush:
				ticker.Stop()
				return
			}
		}
	}()

	return f
}

var _ Logger = (*FileLogger)(nil)

func (f *FileLogger) Close() {
	f.stopAutoflush <- true
	f.writer.Flush()
	f.writer = nil
	f.logFile.Close()
	f.logFile = nil
}

func (f *FileLogger) Debug(a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[DEBUG] %s  ", time.Now().Format(fldf))
	fmt.Fprintln(f.writer, a...)
}

func (f *FileLogger) DebugF(format string, a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[DEBUG] %s  ", time.Now().Format(fldf))
	fmt.Fprintf(f.writer, format, a...)
	fmt.Fprint(f.writer, "\n")
}

func (f *FileLogger) Error(filename string, line int, a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[ERR]   %s  [%s:%d] ", time.Now().Format(fldf), filename, line)
	fmt.Fprintln(f.writer, a...)
}

func (f *FileLogger) ErrorF(filename string, line int, format string, a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[ERR]   %s  [%s:%d] ", time.Now().Format(fldf), filename, line)
	fmt.Fprintf(f.writer, format, a...)
	fmt.Fprint(f.writer, "\n")
}

func (f *FileLogger) Info(a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[INFO]  %s  ", time.Now().Format(fldf))
	fmt.Fprintln(f.writer, a...)
}

func (f *FileLogger) InfoF(format string, a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[INFO]  %s  ", time.Now().Format(fldf))
	fmt.Fprintf(f.writer, format, a...)
	fmt.Fprint(f.writer, "\n")
}

func (f *FileLogger) Warning(a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[WARN]  %s  ", time.Now().Format(fldf))
	fmt.Fprintln(f.writer, a...)
}

func (f *FileLogger) WarningF(format string, a ...any) {
	f.changeLogfileIfNeeded()

	fmt.Fprintf(f.writer, "[WARN]  %s  ", time.Now().Format(fldf))
	fmt.Fprintf(f.writer, format, a...)
	fmt.Fprint(f.writer, "\n")
}

func (f *FileLogger) createCurrentTime() int {
	now := time.Now()
	y, M, d := now.Date()

	return y<<16 | int(M)<<8 | d
}

func (f *FileLogger) changeLogfileIfNeeded() {
	currentTime := f.createCurrentTime()
	if currentTime == f.currentTime {
		return
	}

	filename := filepath.Join(baseDir, fmt.Sprintf("logs_%s.txt", time.Now().Format(flfdf)))
	f.currentTime = f.createCurrentTime()

	file, err := os.Open(filename)
	if err != nil {
		file, err = os.Create(filename)
		if err != nil {
			panic("Couldn't create log-file: " + err.Error())
		}
	}

	f.writer.Flush()
	f.logFile.Close()

	f.logFile = file
	f.writer = bufio.NewWriter(file)
}
