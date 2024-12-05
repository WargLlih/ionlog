package ioncore

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

type controlFlow struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type ionLogger struct {
	controlFlow
	logEngine     *slog.Logger
	writerHandler ionWriter
	reports       chan ionReport
}

type IIonLogger interface {
	LogEngine() *slog.Logger
	SetLogEngine(handler *slog.Logger)

	Targets() []io.Writer
	SetTargets(targets ...io.Writer)

	Stop()
	CreateDefaultLogHandler() slog.Handler
	SendReport(r ionReport)
}

const timeout = 10 * time.Millisecond

var logger *ionLogger

func init() {
	logger = newLogger()

	// using internaly
	slog.SetDefault(slog.New(slog.NewJSONHandler(DefaultOutput, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func newLogger() *ionLogger {
	l := &ionLogger{}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.reports = make(chan ionReport)
	l.logEngine = slog.New(l.CreateDefaultLogHandler())

	go l.handleIonReports()
	return l
}

// Logger returns the logger instance
func Logger() IIonLogger {
	return logger
}

func (i *ionLogger) LogEngine() *slog.Logger {
	return i.logEngine
}

func (i *ionLogger) SetLogEngine(handler *slog.Logger) {
	i.logEngine = handler
}

func (i *ionLogger) Targets() []io.Writer {
	return i.writerHandler.writeTargets
}

func (i *ionLogger) SetTargets(targets ...io.Writer) {
	i.writerHandler.SetTargets(targets...)
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (i *ionLogger) Stop() {
	slog.Debug("Logger stopped (sync triggered)")
	i.cancel()
	i.wg.Wait()
}

// CreateDefaultLogHandler creates a default log handler for the logger
func (i *ionLogger) CreateDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&i.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

// SendReport sends the report to the Logger handler, it will wait for 10ms before returning.
func (i *ionLogger) SendReport(r ionReport) {
	select {
	case <-time.After(timeout):
		slog.Warn(fmt.Sprintf("Failed to send the report (timeout=%v): %v", timeout, r))
		return
	case i.reports <- r:
	}
}

func (i *ionLogger) handleIonReports() {
	i.wg.Add(1)
	defer i.wg.Done()

	for {
		select {
		case <-i.ctx.Done():
			slog.Debug("context done, flushing logs")
			return

		case r := <-i.reports:
			i.log(r.level, r.msg, r.args...)
		}
	}
}

func (i *ionLogger) log(level slog.Level, msg string, args ...any) {
	switch level {
	case slog.LevelDebug:
		i.logEngine.Debug(msg, args...)
	case slog.LevelInfo:
		i.logEngine.Info(msg, args...)
	case slog.LevelWarn:
		i.logEngine.Warn(msg, args...)
	case slog.LevelError:
		i.logEngine.Error(msg, args...)

	default:
		slog.Warn("Unknown log level")
	}
}
