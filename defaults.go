package ionlog

import (
	"io"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
)

// DefaultOutput returns the standard output (stdout)
func DefaultOutput() io.Writer {
	return ioncore.DefaultOutput
}
