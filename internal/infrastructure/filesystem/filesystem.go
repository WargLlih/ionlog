package filesystem

import "os"

type Filesystem struct {
	Stat       func(string) (os.FileInfo, error)
	Mkdir      func(string, os.FileMode) error
	ReadDir    func(string) ([]os.DirEntry, error)
	IsNotExist func(error) bool
	OpenFile   func(name string, flag int, perm os.FileMode) (*os.File, error)
	RemoveFile func(name string) error
}

func NewFileSystem(
	stat func(string) (os.FileInfo, error),
	mkdir func(string, os.FileMode) error,
	readDir func(string) ([]os.DirEntry, error),
	isNotExist func(error) bool,
	openFile func(name string, flag int, perm os.FileMode) (*os.File, error),
	removeFile func(name string) error,
) Filesystem {
	return Filesystem{
		Stat:       stat,
		Mkdir:      mkdir,
		ReadDir:    readDir,
		IsNotExist: isNotExist,
		OpenFile:   openFile,
		RemoveFile: removeFile,
	}
}
