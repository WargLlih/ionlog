package logrotation

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/interfaces"
)

// Start starts the log file rotation service. It blocks until the service is stopped.
func (l *logFileRotation) Start() error {
	l.serviceStatus = interfaces.Running
	defer func() { l.serviceStatus = interfaces.Stopped }()

	if err := validateRotation(l.rotation); err != nil {
		return err
	}

	f, err := l.getActualFile()
	l.logFile = f
	defer l.closeFile()

	l.UnblockWrite()

	if err != nil {
		return err
	}

	// check folder size
	go l.checkFolderSize()

	// every ticker check if the log file needs to be rotated
	var ticker *time.Ticker
	switch l.rotation {
	case Daily:
		ticker = time.NewTicker(8 * time.Hour) // every 8 hours
	case Weekly:
		ticker = time.NewTicker(3 * 24 * time.Hour) // every 3 days
	case Monthly:
		ticker = time.NewTicker(7 * 24 * time.Hour) // every 7 days
	default:
		slog.Error(fmt.Sprintf("rotation was validated but it's invalid: %v", l.rotation))
		return ErrInvalidRotation
	}

	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			slog.Debug("logfile system stopped by context")
			return nil

		case <-ticker.C:
			err := func() error {
				l.BlockWrite()
				defer l.UnblockWrite()

				if err := l.logFile.Close(); err != nil {
					slog.Warn(err.Error())
				}

				f, err := l.getActualFile()
				l.logFile = f

				if err != nil {
					return err
				}

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}
}

// Stop stops the log file rotation service.
func (l *logFileRotation) Stop() {
	l.BlockWrite()
	defer l.UnblockWrite()

	l.cancel()
}

// Status returns the status of the log file rotation service.
func (l *logFileRotation) Status() interfaces.ServiceStatus {
	return l.serviceStatus
}
