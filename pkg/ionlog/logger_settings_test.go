package ionlog

import (
	"io"
	"log/slog"
	"os"
	"testing"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
)

func TestSetLogAttributes(t *testing.T) {
	t.Run("WithTargets", func(t *testing.T) {
		ioncore.Logger().SetLogEngine(slog.New(ioncore.Logger().CreateDefaultLogHandler()))

		w := []io.Writer{ioncore.DefaultOutput, os.Stderr}

		SetLogAttributes(WithTargets(w...))

		for _, writer := range ioncore.Logger().Targets() {
			if writer == nil {
				t.Errorf("Expected writer to be not nil")
			}
		}
	})
}
