// Package logrotaion provides a log file rotation service.
// It is responsible for managing the log file, creating new log files, and rotating the log files.
package logrotation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/filesystem"
	"github.com/IonicHealthUsa/ionlog/internal/interfaces"
)

const logFilePattern = "logfile-%s.log"

type logFileRotation struct {
	filesystem.Filesystem
	logFile       io.WriteCloser
	writeMutex    sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
	serviceStatus interfaces.ServiceStatus
	folder        string
	rotation      PeriodicRotation
}

type ILogFileRotation interface {
	interfaces.IService
	io.Writer
	BlockWrite()
	UnblockWrite()
}

var _ ILogFileRotation = &logFileRotation{}

func NewLogFileRotation(folder string, rotation PeriodicRotation) *logFileRotation {
	c, cancel := context.WithCancel(context.Background())
	return &logFileRotation{
		Filesystem: filesystem.NewFileSystem(
			os.Stat,
			os.Mkdir,
			os.ReadDir,
			os.IsNotExist,
			os.OpenFile,
		),
		ctx:      c,
		cancel:   cancel,
		folder:   folder,
		rotation: rotation,
	}
}

// Write writes the log message to the log file.
func (l *logFileRotation) Write(p []byte) (n int, err error) {
	l.BlockWrite()
	defer l.UnblockWrite()

	if l.logFile == nil {
		return 0, ErrCouldNotGetActualFile
	}

	return l.logFile.Write(p)
}

func (l *logFileRotation) BlockWrite() {
	l.writeMutex.Lock()
}

func (l *logFileRotation) UnblockWrite() {
	l.writeMutex.Unlock()
}

// closeFile closes the log file.
func (l *logFileRotation) closeFile() {
	if l.logFile != nil {
		if err := l.logFile.Close(); err != nil {
			slog.Warn(err.Error())
		}
		l.logFile = nil
	}
}

// getActualFile gets the actual log file to write to.
// It returns the actual file and an error if it couldn't get the actual file.
// It checks if the folder exists, if it doesn't exist it creates the folder.
// It checks if the folder has any log files, if it doesn't have any log files it creates a new log file.
// It checks if the log file needs to be rotated based on the rotation type.
func (l *logFileRotation) getActualFile() (*os.File, error) {
	var actualFile *os.File

	exists, err := l.checkFolder(l.folder)
	if err != nil {
		return nil, err
	}

	// folder exists
	if exists {
		files, err := l.getAllfiles(l.folder)
		if err != nil {
			return nil, err
		}

		mostRecent, err := getMostRecentLogFile(files)

		if errors.Is(err, ErrNoLogFileFound) {
			// no file in the folder, create a new file.
			newFilename := createNewLogFilename()
			actualFile, err = l.createFileInFolder(newFilename)
			if err != nil {
				slog.Debug(err.Error())
				return nil, ErrFailedToCreateFile
			}
			return actualFile, nil
		}

		if err != nil {
			slog.Debug(err.Error())
			return nil, ErrCouldNotGetActualFile
		}

		// File exists, check if it needs to be rotated
		fileDate, err := getFileDate(mostRecent)
		if err != nil {
			slog.Debug(err.Error())
			return nil, ErrCouldNotGetActualFile
		}

		if checkRotation(l.rotation, fileDate) {
			newFilename := createNewLogFilename()
			actualFile, err = l.createFileInFolder(newFilename)
			if err != nil {
				slog.Debug(err.Error())
				return nil, ErrFailedToCreateFile
			}
			return actualFile, nil
		}

		// File doesn't need to be rotated, open the file
		actualFile, err = l.OpenFile(
			filepath.Join(l.folder, mostRecent),
			os.O_WRONLY|os.O_APPEND,
			0644,
		)

		if err != nil {
			slog.Debug(err.Error())
			return nil, ErrCouldNotGetActualFile
		}
		return actualFile, nil

	} else {
		if err := l.Mkdir(l.folder, 0755); err != nil {
			slog.Debug(err.Error())
			return nil, ErrCouldNotCreateFolder
		}

		newFilename := createNewLogFilename()
		actualFile, err = l.createFileInFolder(newFilename)
		if err != nil {
			slog.Debug(err.Error())
			return nil, ErrFailedToCreateFile
		}
		return actualFile, nil
	}
}

// checkFolder checks if the folder exists, it returns true if the folder exists
// false if it doesn't exist and an error if it couldn't check the folder status.
func (l *logFileRotation) checkFolder(folder string) (bool, error) {
	// Check if folder exists
	_, err := l.Stat(folder)

	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	slog.Debug(err.Error())
	return false, ErrCouldNotCheckFolderStatus
}

// getAllfiles gets all the files in the folder.
// It returns a list of filenames and an error if it couldn't read the folder.
func (l *logFileRotation) getAllfiles(folder string) ([]string, error) {
	files, err := l.ReadDir(folder)
	if err != nil {
		slog.Debug(err.Error())
		return nil, ErrFailedToReadFolder
	}

	filePattern := regexp.MustCompile(`logfile-\d{4}-\d{2}-\d{2}.log`)

	var filenames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !filePattern.MatchString(file.Name()) {
			slog.Warn(fmt.Sprintf("file: %s is not a valid log file. Skipping.", file.Name()))
			continue
		}

		filenames = append(filenames, file.Name())
	}

	return filenames, nil
}

// createFileInFolder creates a new log file in the specified folder.
func (l *logFileRotation) createFileInFolder(filename string) (*os.File, error) {
	filePath := filepath.Join(l.folder, filename)
	return l.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
}

// getMostRecentLogFile gets the most recent log file from the list of files.
// It returns the most recent log filename and an error if no log file was found.
func getMostRecentLogFile(files []string) (string, error) {
	var mostRecent string
	var latestTime time.Time

	for _, file := range files {
		fileTime, err := getFileDate(file)
		if err != nil {
			slog.Warn(fmt.Sprintf("Failed to get file date for file: %s. Skipping.", file))
			continue
		}

		if fileTime.After(latestTime) {
			latestTime = fileTime
			mostRecent = file
		}
	}

	if mostRecent == "" {
		return "", ErrNoLogFileFound
	}

	return mostRecent, nil
}

// getFileDate gets the date from the log file name.
// It returns the date and an error if the date couldn't be parsed.
func getFileDate(file string) (time.Time, error) {
	dateStr := strings.TrimPrefix(strings.TrimSuffix(file, ".log"), "logfile-")
	return time.Parse(time.DateOnly, dateStr)
}

// createNewLogFilename creates a new log file name based on the current date.
func createNewLogFilename() string {
	return fmt.Sprintf(logFilePattern, time.Now().Format(time.DateOnly))
}

// checkRotation checks if the log file needs to be rotated based on the rotation type.
// It returns true if the log file needs to be rotated and false if it doesn't.
func checkRotation(rotation PeriodicRotation, fileDate time.Time) bool {
	now := time.Now()
	y, m, d := now.Date()
	_, w := now.UTC().ISOWeek()

	switch rotation {
	case Daily:
		return fileDate.Year() != y || fileDate.Month() != m || fileDate.Day() != d
	case Weekly:
		// Warning: Read the ISOWeek() documentation to check any inconsistencies
		// in the end of a year or the beginning of a new year.
		_, fileW := fileDate.UTC().ISOWeek()
		return fileDate.Year() != y || fileW != w
	case Monthly:
		return fileDate.Year() != y || fileDate.Month() != m
	default:
		slog.Error(fmt.Sprintf("rotation was validated but it's invalid: %v", rotation))
		return false
	}
}

func validateRotation(rotation PeriodicRotation) error {
	switch rotation {
	case Daily, Weekly, Monthly:
		return nil
	default:
		slog.Debug(fmt.Sprintf("Invalid rotation type: %v", rotation))
		return ErrInvalidRotation
	}
}
