package logcore

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	_logger := newLogger()

	if _logger == nil {
		t.Errorf("Expected logger to be not nil")
	}
}

func TestCreateDefaultLogHandler(t *testing.T) {
	_logger := newLogger()
	handler := _logger.CreateDefaultLogHandler()

	if handler == nil {
		t.Errorf("Expected handler to be not nil")
	}
}
