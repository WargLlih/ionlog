package ionlog

import (
	"fmt"
	"io"
	"os"
)

type ionWriter struct {
	writeTargets []io.Writer
}

// Write writes the contents of p to all writeTargets
// if any writeTarget returns an error, the error is returned with the number of bytes written
// if no writeTarget returns an error, this function returns no err = nil and n = 0
func (w *ionWriter) Write(p []byte) (n int, err error) {
	for i, target := range w.writeTargets {

		if target == nil {
			ionInternalLogger.Error("Expected writer to be not nil")
			continue
		}

		// It will save the latest failure error while continue writing to other writeTargets
		// latter, it will return the latest failure error
		_n, _err := target.Write(p)
		if _err != nil {
			n = _n
			err = _err
			ionInternalLogger.Error(fmt.Sprintf("Failed to write to in the %v° target, error: %v", i+1, err))
		}
	}
	return
}

func DefaultOutput() io.Writer {
	return os.Stdout
}
