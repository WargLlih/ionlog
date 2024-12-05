package benchmark_test

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog/pkg/ionlog"
)

type fakkeWriter struct {
}

func (f *fakkeWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func BenchmarkIonlog(b *testing.B) {

	ionlog.SetLogAttributes(
		ionlog.WithTargets(&fakkeWriter{}),
		ionlog.WithStaicFields(map[string]string{
			"computer-id": "1234",
		}),
	)

	for i := 0; i < b.N; i++ {
		ionlog.Info("This log level is: %v", "info")
	}
}
