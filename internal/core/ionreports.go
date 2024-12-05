package ioncore

import (
	"log/slog"
)

type ionReport struct {
	level slog.Level
	msg   string
	args  []any
}

func NewIonReport(level slog.Level, msg string, args []any) ionReport {
	return ionReport{
		level: level,
		msg:   msg,
		args:  args,
	}
}
