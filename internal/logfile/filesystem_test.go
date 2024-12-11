package ionlogfile

import (
	"os"
	"testing"
)

func TestNewFileSystem(t *testing.T) {
	t.Run("New File System", func(t *testing.T) {
		fs := newFileSystem(os.Stat, os.Mkdir, os.ReadDir, os.IsNotExist, os.OpenFile)
		if fs.stat == nil {
			t.Errorf("stat function is nil")
		}
		if fs.mkdir == nil {
			t.Errorf("mkdir function is nil")
		}
		if fs.readDir == nil {
			t.Errorf("readDir function is nil")
		}
		if fs.isNotExist == nil {
			t.Errorf("isNotExist function is nil")
		}
		if fs.openFile == nil {
			t.Errorf("openFile function is nil")
		}
	})
}
