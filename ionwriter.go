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
// This function returns no error nor the number of bytes written
func (w *ionWriter) Write(p []byte) (int, error) {
	for i, target := range w.writeTargets {
		if target == nil {
			ionInternalLogger.Error(fmt.Sprintf("Expected the %v° target to be not nil", i+1))
			continue
		}

		_, err := target.Write(p)
		if err != nil {
			ionInternalLogger.Error(fmt.Sprintf("Failed to write to in the %v° target, error: %v", i+1, err))
		}
	}

	return 0, nil
}

func DefaultOutput() io.Writer {
	return os.Stdout
}
