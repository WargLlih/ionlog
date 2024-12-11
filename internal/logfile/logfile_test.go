package ionlogfile

import (
	"io"
	"testing"
)

func TestNewLogFileRotation(t *testing.T) {
	t.Run("New Log File Rotation", func(t *testing.T) {
		lfr := NewLogFileRotation("folder", Daily)
		if lfr == nil {
			t.Errorf("NewLogFileRotation() failed")
		}

		if lfr.filesystem.stat == nil {
			t.Errorf("stat function is nil")
		}
		if lfr.filesystem.mkdir == nil {
			t.Errorf("mkdir function is nil")
		}
		if lfr.filesystem.readDir == nil {
			t.Errorf("readDir function is nil")
		}
		if lfr.filesystem.isNotExist == nil {
			t.Errorf("isNotExist function is nil")
		}
		if lfr.filesystem.openFile == nil {
			t.Errorf("openFile function is nil")
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
			t.Errorf("NewLogFileRotation() failed")
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
			t.Errorf("NewLogFileRotation() failed")
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
			t.Errorf("NewLogFileRotation() failed")
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

