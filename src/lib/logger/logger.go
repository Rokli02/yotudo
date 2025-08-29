package logger

import "runtime"

const (
	Console_Type  = "console"
	Database_Type = "db"
	File_Type     = "file"
)

type Logger interface {
	Info(a ...any)
	InfoF(format string, a ...any)
	Debug(a ...any)
	DebugF(format string, a ...any)
	WarningF(format string, a ...any)
	Warning(a ...any)
	ErrorF(filename string, line int, format string, a ...any)
	Error(filename string, line int, a ...any)
	Close()
}

var loggers []Logger
var logLevel uint8

const (
	error_log_level uint8 = iota
	warning_log_level
	info_log_level
	debug_log_level
)

func InitializeLogger(level string, types []string) (chosenLoggers []Logger, closeLoggers func()) {
	switch level {
	case "debug":
		logLevel = debug_log_level
	case "info":
		logLevel = info_log_level
	case "error":
		logLevel = error_log_level
	case "warning":
		fallthrough
	default:
		logLevel = warning_log_level
	}

	loggers = make([]Logger, 0)
	for _, t := range types {
		switch t {
		case Console_Type:
			loggers = append(loggers, NewConsoleLogger())
		case Database_Type:
			panic("NOT_IMPLEMENTED_DATABASE_LOGGER")
		case File_Type:
			panic("NOT_IMPLEMENTED_FILE_LOGGER")
		}
	}
	chosenLoggers = loggers
	closeLoggers = func() {
		for _, logger := range loggers {
			logger.Close()
		}

		loggers = nil
	}

	return
}

func Info(a ...any) {
	if logLevel < info_log_level {
		return
	}

	for _, logger := range loggers {
		logger.Info(a...)
	}
}

func InfoF(format string, a ...any) {
	if logLevel < info_log_level {
		return
	}

	for _, logger := range loggers {
		logger.InfoF(format, a...)
	}
}

func Debug(a ...any) {
	if logLevel < debug_log_level {
		return
	}

	for _, logger := range loggers {
		logger.Debug(a...)
	}
}

func DebugF(format string, a ...any) {
	if logLevel < debug_log_level {
		return
	}

	for _, logger := range loggers {
		logger.DebugF(format, a...)
	}
}

func WarningF(format string, a ...any) {
	if logLevel < warning_log_level {
		return
	}

	for _, logger := range loggers {
		logger.WarningF(format, a...)
	}
}

func Warning(a ...any) {
	if logLevel < warning_log_level {
		return
	}

	for _, logger := range loggers {
		logger.Warning(a...)
	}
}

func ErrorF(format string, a ...any) {
	// Unsigned integer can't be negative and error has the value 0 -> check is not needed here
	_, filename, line, _ := runtime.Caller(1)

	for _, logger := range loggers {
		logger.ErrorF(filename, line, format, a...)
	}
}

func Error(a ...any) {
	// Unsigned integer can't be negative and error has the value 0 -> check is not needed here
	_, filename, line, _ := runtime.Caller(1)

	for _, logger := range loggers {
		logger.Error(filename, line, a...)
	}
}
