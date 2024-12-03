package ionlog

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestSetLogAttributes(t *testing.T) {
	t.Run("WithTargets", func(t *testing.T) {
		logger = newLogger()
		logger.logHandler = slog.New(logger.createDefaultLogHandler())

		w := []io.Writer{Stdout(), os.Stderr}

		SetLogAttributes(WithTargets(w...))

		for _, writer := range logger.writerHandler.writeTargets {
			if writer == nil {
				t.Errorf("Expected writer to be not nil")
			}
		}
	})
}
