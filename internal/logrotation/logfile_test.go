package logrotation

import (
	"io"
	"testing"
)

func TestNewLogFileRotation(t *testing.T) {
	t.Run("New Log File Rotation", func(t *testing.T) {
		lfr := NewLogFileRotation("folder", Daily)
		if lfr == nil {
			t.Fatalf("NewLogFileRotation() failed")
		}

		if lfr.Filesystem.Stat == nil {
			t.Errorf("Stat function is nil")
		}
		if lfr.Filesystem.Mkdir == nil {
			t.Errorf("Mkdir function is nil")
		}
		if lfr.Filesystem.ReadDir == nil {
			t.Errorf("ReadDir function is nil")
		}
		if lfr.Filesystem.IsNotExist == nil {
			t.Errorf("IsNotExist function is nil")
		}
		if lfr.Filesystem.OpenFile == nil {
			t.Errorf("OpenFile function is nil")
		}

		if lfr.ctx == nil {
			t.Errorf("ctx is nil")
		}
		if lfr.cancel == nil {
			t.Errorf("cancel is nil")
		}

		if lfr.folder != "folder" {
			t.Errorf("folder is not set")
		}

		if lfr.rotation != Daily {
			t.Errorf("rotation is not set")
		}
	})
}

type mockFileSystem struct {
	writeError error
	closeError error
}

var _ io.WriteCloser = &mockFileSystem{}

func (mfs *mockFileSystem) Write(p []byte) (n int, err error) {
	return 0, mfs.writeError
}

func (mfs *mockFileSystem) Close() error {
	return mfs.closeError
}

func TestLogFileRotation_Write(t *testing.T) {
	t.Run("Write No File", func(t *testing.T) {
		lfr := NewLogFileRotation("folder", Daily)
		if lfr == nil {
			t.Fatalf("NewLogFileRotation() failed")
		}

		lfr.logFile = nil

		_, err := lfr.Write([]byte("test"))
		if err != ErrCouldNotGetActualFile {
			t.Errorf("Expected error: %v, got: %v", ErrCouldNotGetActualFile, err)
		}
	})

	t.Run("Write Success", func(t *testing.T) {
		lfr := NewLogFileRotation("folder", Daily)
		if lfr == nil {
			t.Fatalf("NewLogFileRotation() failed")
		}

		lfr.logFile = &mockFileSystem{}

		_, err := lfr.Write([]byte("test"))
		if err != nil {
			t.Errorf("Expected error: %v, got: %v", nil, err)
		}
	})

	t.Run("Write Write Error", func(t *testing.T) {
		lfr := NewLogFileRotation("folder", Daily)
		if lfr == nil {
			t.Fatalf("NewLogFileRotation() failed")
		}

		lfr.logFile = &mockFileSystem{
			writeError: io.ErrShortWrite,
		}

		_, err := lfr.Write([]byte("test"))
		if err == nil {
			t.Errorf("Expected error: %v, got: %v", nil, err)
		}
	})
}
