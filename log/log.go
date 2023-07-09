package log

import (
	golog "log"
)

//go: cover ignore
func Info(format string, v ...any) {
	golog.Print(msgFormat(msgInfo, format, v...))
}

//go: cover ignore
func Warning(format string, v ...any) {
	golog.Print(msgFormat(msgWarning, format, v...))
}

//go: cover ignore
func Debug(format string, v ...any) {
	golog.Print(msgFormat(msgDebug, format, v...))
}

//go: cover ignore
func Notice(format string, v ...any) {
	golog.Print(msgFormat(msgNotice, format, v...))
}

//go: cover ignore
func Error(format string, v ...any) {
	golog.Print(msgFormat(msgError, format, v...))
}

//go: cover ignore
func Emergency(format string, v ...any) {
	golog.Print(msgFormat(msgEmergency, format, v...))
}

//go: cover ignore
func Critical(format string, v ...any) {
	golog.Print(msgFormat(msgCritical, format, v...))
}

//go: cover ignore
func Alert(format string, v ...any) {
	golog.Print(msgFormat(msgAlert, format, v...))
}
