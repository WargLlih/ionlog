package benchmark

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog/pkg/ionlog"
)

type fakkeWriter struct{}

func (fakkeWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func BenchmarkSimpleLogs(b *testing.B) {
	// Set the log attributes, and other configurations
	ionlog.SetLogAttributes(
		ionlog.WithTargets(fakkeWriter{}),
	)

	// Start the logger service
	ionlog.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ionlog.Infof("This log level is: %v", "info")
	}
}
