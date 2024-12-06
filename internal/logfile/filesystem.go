package ionlogfile

import "os"

type filesystem struct {
	stat       func(string) (os.FileInfo, error)
	mkdir      func(string, os.FileMode) error
	readDir    func(string) ([]os.DirEntry, error)
	isNotExist func(error) bool
	openFile   func(name string, flag int, perm os.FileMode) (*os.File, error)
}

func newFileSystem(
	stat func(string) (os.FileInfo, error),
	mkdir func(string, os.FileMode) error,
	readDir func(string) ([]os.DirEntry, error),
	isNotExist func(error) bool,
	openFile func(name string, flag int, perm os.FileMode) (*os.File, error),
) filesystem {
	return filesystem{
		stat:       stat,
		mkdir:      mkdir,
		readDir:    readDir,
		isNotExist: isNotExist,
		openFile:   openFile,
	}
}
