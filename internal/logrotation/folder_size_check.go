package logrotation

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

func (l *logFileRotation) checkFolderSize() {
	routine := func() {
		l.BlockWrite()
		defer l.UnblockWrite()

		size, err := l.getFolderSize()
		if err != nil {
			slog.Error(err.Error())
		}
		if size > l.maxFolderSize {
			files, err := l.getAllfiles(l.folder)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			oldestFile, err := getOldestLogFile(files)
			if err != nil {
				slog.Error(err.Error())
				return
			}
			slog.Info(fmt.Sprintf("Removing file: %s", oldestFile))
			if err = l.RemoveFile(filepath.Join(l.folder, oldestFile)); err != nil {
				slog.Error(err.Error())
				return
			}

			// check if it need to create a new file
			actualFile, err := l.getActualFile()
			if err != nil {
				slog.Error(err.Error())
				return
			}
			l.logFile = actualFile
		}
	}
	for {
		select {
		case <-l.ctx.Done():
			slog.Debug("folder size check system stopped by context")
			return

		case <-time.After(30 * time.Minute):
			routine()
		}
	}
}

func getOldestLogFile(files []string) (string, error) {
	var oldestFile string
	var oldestTime = time.Now()

	for _, file := range files {
		fileTime, err := getFileDate(file)
		if err != nil {
			slog.Warn(fmt.Sprintf("Failed to get file date for file: %s. Skipping.", file))
			continue
		}

		fmt.Printf("fileTime: %v | oldestTime: %v\n", fileTime, oldestTime)

		if fileTime.Before(oldestTime) {
			oldestTime = fileTime
			oldestFile = file
		}
	}

	if oldestFile == "" {
		return "", ErrNoLogFileFound
	}

	return oldestFile, nil
}

func (l *logFileRotation) getFolderSize() (uint, error) {
	var size int64
	err := filepath.Walk(l.folder, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return uint(size), err
}
