package ioncore

import (
	"log/slog"

	ionlogfile "github.com/IonicHealthUsa/ionlog/internal/logfile"
	ionservice "github.com/IonicHealthUsa/ionlog/internal/service"
)

// Start starts the logger service, it blocks until the service is stopped
func (i *ionLogger) Start() error {
	i.serviceStatus = ionservice.Running
	defer func() { i.serviceStatus = ionservice.Stopped }()

	go func() {
		if i.rotationPeriod == ionlogfile.NoAutoRotate {
			return
		}

		logRotate := ionlogfile.NewLogFileRotation(i.folder, i.rotationPeriod)
		if err := logRotate.Start(); err != nil {
			slog.Error(err.Error())
			return
		}
	}()

	i.handleIonReports()
	return nil
}

// Status returns the status of the logger service
func (i *ionLogger) Status() ionservice.ServiceStatus {
	return i.serviceStatus
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (i *ionLogger) Stop() {
	slog.Debug("Logger stopped (sync triggered)")
	i.cancel()
	i.wg.Wait()
}
