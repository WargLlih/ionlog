# IonLog

## Usage
```go
package main

import (
	"github.com/IonicHealthUsa/ionlog/pkg/ionlog"
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

		// (Optional) WithStaicFields sets the static fields for the logger, every log will have these fields.
		ionlog.WithStaicFields(map[string]string{
			"computer-id": "1234",
			// your custom fields
		}),

		// (Optional) WithLogFileRotation sets the log file rotation period and the folder where the log files will be stored.
		// This is a internal log file rotation system, when optionally used, it will append the log file to the targets, and
		// will rotate it automatically.
		ionlog.WithLogFileRotation("logs", ionlog.Daily),
	)

	// Start the logger service
	ionlog.Start()

	// Stops the logger service when the main function ends
	defer ionlog.Stop()

	// output: {"time":"2024-12-06T20:59:47.252944832-03:00","level":"INFO","msg":"This log level is: info","computer-id":"1234","package":"main","function":"main","file":"main.go","line":38}
	ionlog.Info("This log level is: %v", "info")
	ionlog.Error("This log level is: %v", "error")
	ionlog.Warn("This log level is: %v", "warn")
	ionlog.Debug("This log level is: %v", "debug")
}
```
