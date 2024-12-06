package ionlogfile

import "errors"

var (
	ErrNoLogFileFound            = errors.New("no valid log files found")
	ErrCouldNotCreateFolder      = errors.New("could not create folder")
	ErrCouldNotCheckFolderStatus = errors.New("could not check folder status")
	ErrFailedToReadFolder        = errors.New("failed to read folder")
	ErrFailedToCreateFile        = errors.New("failed to create file")
	ErrCouldNotGetActualFile     = errors.New("could not get actual file")
	ErrInvalidRotation           = errors.New("invalid rotation")
)
