package benchmark

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

func BenchmarkSimpleLogs(b *testing.B) {
	// Set the log attributes, and other configurations
	ionlog.SetLogAttributes(
		ionlog.WithTargets(ionlog.DefaultOutput()),
	)

	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ionlog.Infof("This log test: %v", i)
	}
}
