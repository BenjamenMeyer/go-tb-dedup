package log

import (
	"fmt"
)

const (
	msgInfo      = "INFO"
	msgWarning   = "WARNING"
	msgDebug     = "DEBUG"
	msgNotice    = "NOTICE"
	msgError     = "ERROR"
	msgEmergency = "EMERGENCY"
	msgCritical  = "CRITICAL"
	msgAlert     = "ALERT"
)

func msgFormat(level string, format string, v ...any) string {
	return fmt.Sprintf("%s: %s", level, fmt.Sprintf(format, v...))
}

func msgUnitFormat(unit string, level string, format string, v ...any) string {
	return fmt.Sprintf("%s(%s): %s", level, unit, fmt.Sprintf(format, v...))
}
