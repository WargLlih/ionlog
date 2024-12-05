# IonLog

## Usage
```go
package main

import (
	"github.com/IonicHealthUsa/ionlog/pkg/ionlog"
)

func main() {
	defer ionlog.Stop()

	ionlog.SetLogAttributes(
		ionlog.WithTargets(ionlog.DefaultOutput() /* socket, files ... */),
		ionlog.WithStaicFields(map[string]string{
			"computer-id": "1234",
			"...":         "...",
		}),
	)

	ionlog.Debug("This log level is: %v", "debug")
	ionlog.Info("This log level is: %v", "info")
	ionlog.Error("This log level is: %v", "error")
	ionlog.Warn("This log level is: %v", "warn")
}
```
