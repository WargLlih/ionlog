package ionlog

import (
	"fmt"
	"log/slog"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
)

func Start() {
	ioncore.Logger().Start()
}

func Stop() {
	ioncore.Logger().Stop()
}

// Info logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Info(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Error logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Error(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelError, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Warn logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warn(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelWarn, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Debug logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debug(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelDebug, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}
