package ionlog

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type ionReports struct {
	level slog.Level
	msg   string
	args  []any
}

type ionLogger struct {
	logHandler    *slog.Logger
	writerHandler ionWriter
	ctx           context.Context
	reports       chan ionReports
}

var logger *ionLogger = newLogger()
var ionInternalLogger = slog.New(slog.NewJSONHandler(Stdout(), &slog.HandlerOptions{Level: slog.LevelDebug}))

func newLogger() *ionLogger {
	l := &ionLogger{}

	l.reports = make(chan ionReports)
	l.logHandler = slog.New(l.createDefaultLogHandler())

	return l
}

// Init initializes the logger handler.
// ctx is the context that will be used to stop the logger handler.
func Init(ctx context.Context) {
	logger.ctx = ctx
	logger.handleIonReports()
}

func (l *ionLogger) createDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&l.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

func (i *ionLogger) handleIonReports() {
	for {
		select {
		case <-i.ctx.Done():
			ionInternalLogger.Info("Logger has been stopped (context is done)")
			return

		case r := <-i.reports:
			switch r.level {
			case slog.LevelInfo:
				i.logHandler.Info(r.msg, r.args...)
			case slog.LevelError:
				i.logHandler.Error(r.msg, r.args...)
			case slog.LevelWarn:
				i.logHandler.Warn(r.msg, r.args...)
			case slog.LevelDebug:
				i.logHandler.Debug(r.msg, r.args...)
			}
		}
	}
}

func sendReport(r ionReports) {
	select {
	case <-time.After(10 * time.Millisecond):
		ionInternalLogger.Warn(fmt.Sprintf("Failed to send the report (timeout=10ms): %v", r))
		return
	case logger.reports <- r:
	}
}

func Info(msg string, args ...any) {
	sendReport(ionReports{slog.LevelInfo, fmt.Sprintf(msg, args...), getRecordInformation()})
}

func Error(msg string, args ...any) {
	sendReport(ionReports{slog.LevelError, fmt.Sprintf(msg, args...), getRecordInformation()})
}

func Warn(msg string, args ...any) {
	sendReport(ionReports{slog.LevelWarn, fmt.Sprintf(msg, args...), getRecordInformation()})
}

func Debug(msg string, args ...any) {
	sendReport(ionReports{slog.LevelDebug, fmt.Sprintf(msg, args...), getRecordInformation()})
}
