# IonLog

# Usage
```go
package main

import (
	"github.com/IonicHealthUsa/ionlog"
)

func main() {
	// Set the log attributes, and other configurations
	ionlog.SetLogAttributes(
		// WithTargets sets the write targets for the logger, every log will be written
		// to these targets.
		ionlog.WithTargets(
			ionlog.DefaultOutput(),
			// a websocket
			// a file
			// your custom writer
		),

		// (Optional) WithStaticFields sets the static fields for the logger, every log will have these fields.
		ionlog.WithStaticFields(map[string]string{
			"computer-id": "1234",
			// your custom fields
		}),

		// (Optional) WithLogFileRotation enables log file rotation, specifying the directory where log files will be stored, the maximum size of the log folder in bytes, and the rotation frequency.
		// This internal log rotation system appends the log file to the specified targets and automatically rotates logs based on the provided configuration,
		// ensuring the total size of the log folder does not exceed the specified maximum (e.g., 10MB in this case).
		ionlog.WithLogFileRotation("logs", 10*ionlog.Mebibyte, ionlog.Daily),
	)

	// Start the logger service
	ionlog.Start()

	// Stops the logger service when the main function ends
	defer ionlog.Stop()

	// output: {"time":"2024-12-06T20:59:47.252944832-03:00","level":"INFO","msg":"This log level is: info","computer-id":"1234","package":"main","function":"main","file":"main.go","line":38}
	ionlog.Infof("This log level is: %v", "info")
	ionlog.Errorf("This log level is: %v", "error")
	ionlog.Warnf("This log level is: %v", "warn")
	ionlog.Debugf("This log level is: %v", "debug")

	ionlog.Info("This log level is a simple info log")
	ionlog.Error("This log level is a simple error log")
	ionlog.Warn("This log level is a simple warn log")
	ionlog.Debug("This log level is a simple debug log")

	status := "NOT OK"
	for i := 0; i < 10; i++ {
		ionlog.LogOnceInfo("Process Started!")   // This will be logged only once
		ionlog.LogOnChangeDebugf("count: %v", i) // Log every time i changes
		if i == 5 {
			status = "OK"
		}
		ionlog.LogOnChangeInfof("status: %v", status) // Log once "NOT OK", log once "OK"
	}
}
```
# Library Import and Configuration:
The library is imported from github.com/IonicHealthUsa/ionlog/pkg/ionlog.  
Configuration is done using SetLogAttributes() method with several options:

## Log Targets:
ionlog.WithTargets() allows setting multiple log output destinations.  
It supports:

- Default output
- Websocket
- File writing
- Custom writers

## Static Fields:
WithStaticFields() adds consistent metadata to all log entries.  
In this example, a "computer-id" field is added to every log

## Log File Rotation:
WithLogFileRotation() configures automatic log file management.
- Sets log file storage location to "logs" folder
- Rotation period set to daily

# Logging Methods:
Standard log levels:
- Infof(msg, args)
- Errorf(msg, args)
- Warnf(msg, args)
- Debugf(msg, args)
- Info(msg)
- Error(msg)
- Warn(msg)
- Debug(msg)
Special logging methods:  

Logs a message only once:
- LogOnceInfof(msg, args)
- LogOnceErrorf(msg, args)
- LogOnceWarnf(msg, args)
- LogOnceDebugf(msg, args)
- LogOnceInfo(msg)
- LogOnceError(msg)
- LogOnceWarn(msg)
- LogOnceDebug(msg)

Logs when the message changes:
- LogOnChangeInfof(msg, args)
- LogOnChangeErrorf(msg, args)
- LogOnChangeWarnf(msg, args)
- LogOnChangeDebugf(msg, args)
- LogOnChangeInfo(msg)
- LogOnChangeError(msg)
- LogOnChangeWarn(msg)
- LogOnChangeDebug(msg)

# Lifecycle Management:
- Start() initializes the logger
- Stop() (deferred) closes the logger when the program ends

# Log Format:
- Produces JSON-formatted logs
- Includes timestamp, log level, message
- Adds static fields, package, function, file, and line information

# Internal Logging system:
- Internal logs are handled by the slog package, and outputed to the os.Stdout by default.
