package log

import (
	golog "log"
)

type outputInterface func(a ...any) // (n int, err error)

var logOutput outputInterface = golog.Print

func UnitInfo(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgInfo, format, v...))
}

func UnitWarning(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgWarning, format, v...))
}

func UnitDebug(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgDebug, format, v...))
}

func UnitNotice(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgNotice, format, v...))
}

func UnitError(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgError, format, v...))
}

func UnitEmergency(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgEmergency, format, v...))
}

func UnitCritical(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgCritical, format, v...))
}

func UnitAlert(unit string, format string, v ...any) {
	logOutput(msgUnitFormat(unit, msgAlert, format, v...))
}
