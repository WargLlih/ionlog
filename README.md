# ionlog

A flexible and structured logging library for Go with dynamic controls.

## Installation

```bash
go get github.com/IonicHealthUsa/ionlog
```

# Basic Usage
```go
package main

import "github.com/IonicHealthUsa/ionlog"

func main() {
    ionlog.SetLogAttributes(
        ionlog.WithTargets(ionlog.DefaultOutput()), // Log to console
        ionlog.WithStaticFields(map[string]string{"service": "my-app"}),
        ionlog.WithLogFileRotation("logs", 10*ionlog.Mebibyte, ionlog.Daily),
    )

    ionlog.Start()
    defer ionlog.Stop()

    ionlog.Info("Application started")
}
```

# Advanced Usage
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

# Key Features
## Configuration Options

### Targets: Log to multiple destinations (console, files, websockets, custom writers).
```go
ionlog.WithTargets(ionlog.DefaultOutput(), myCustomWriter)
```

### Static Fields: Add fixed fields to all logs (e.g., service name, environment).
```go
ionlog.WithStaticFields(map[string]string{"env": "production"})
```

### Log Rotation: Auto-rotate logs by size and time.
```go
ionlog.WithLogFileRotation("logs", 100*ionlog.Mebibyte, ionlog.Hourly)
```

## Logging Functions
- Levels: Debug, Info, Warn, Error.
```go
ionlog.Infof("User %s logged in", "Alice")
ionlog.Error("Connection failed")
```

## Structured Output: Logs are emitted as JSON with metadata:
```json
{
	"time":"2024-12-06T20:59:47.252944832-03:00",
	"level":"INFO",
	"msg": "User Alice logged in",
	"service-id":"0xcafe",
	"package":"main",
	"function":"main",
	"file":"main.go",
	"line":42
}
```

## Special Logging

### Log Once: Write a message only once during execution.
```go
ionlog.LogOnceInfo("Initialization complete")
```

### Log on Change: Only log when the value changes.
```go
status := "STARTING"
ionlog.LogOnChangeInfof("status: %s", status) // Logs once
ionlog.LogOnChangeInfof("status: %s", status) // Will not log

status = "RUNNING"
ionlog.LogOnChangeInfof("status: %s", status) // Logs again
ionlog.LogOnChangeInfof("status: %s", status) // Will not log
```

## Lifecycle Management:

- Start() initializes the logger
- Stop() closes the logger when the program ends


# Internal Logging system:
- Internal logs are handled by the slog package, and outputed to the os.Stdout by default.
