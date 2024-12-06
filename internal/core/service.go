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

	// user has chosen to auto rotate the log file
	if i.rotationPeriod != ionlogfile.NoAutoRotate {
		i.logRotateService = ionlogfile.NewLogFileRotation(i.folder, i.rotationPeriod)

		// logRotateService is manages a file, so it is a target...
		i.SetTargets(append(i.Targets(), i.logRotateService)...)

		// block until the log rotate service sets up the file to write to.
		i.logRotateService.BlockWrite()
	}

	go func() {
		if i.logRotateService == nil {
			return
		}

		i.wg.Add(1)
		defer i.wg.Done()

		if err := i.logRotateService.Start(); err != nil {
			slog.Error(err.Error())
			return
		}
	}()

	go func() {
		i.wg.Add(1)
		defer i.wg.Done()

		i.handleIonReports()
	}()

	return nil
}

// Status returns the status of the logger service
func (i *ionLogger) Status() ionservice.ServiceStatus {
	return i.serviceStatus
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (i *ionLogger) Stop() {
	slog.Debug("Logger service stopping...")
	i.cancel()

	if i.logRotateService != nil {
		i.logRotateService.Stop()
	}

	i.wg.Wait()
}
