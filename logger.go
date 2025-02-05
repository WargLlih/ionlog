package ionlog

import (
	"fmt"
	"log/slog"

	"github.com/IonicHealthUsa/ionlog/internal/logcore"
	"github.com/IonicHealthUsa/ionlog/internal/usecases"
)

func Start() {
	logcore.Logger().Start()
}

func Stop() {
	logcore.Logger().Stop()
}

// Info logs a message with level info.
func Info(msg any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelInfo, msg.(string), logcore.GetRecordInformation()))
}

// Error logs a message with level error.
func Error(msg any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelError, msg.(string), logcore.GetRecordInformation()))
}

// Warn logs a message with level warn.
func Warn(msg any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelWarn, msg.(string), logcore.GetRecordInformation()))
}

// Debug logs a message with level debug.
func Debug(msg any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelDebug, msg.(string), logcore.GetRecordInformation()))
}

// LogOnceInfo logs a message with level info only once.
func LogOnceInfo(msg any) {
	logOnce(slog.LevelInfo, msg.(string))
}

// LogOnceError logs a message with level info only once.
func LogOnceError(msg any) {
	logOnce(slog.LevelError, msg.(string))
}

// LogOnceWarn logs a message with level warn only once.
func LogOnceWarn(msg any) {
	logOnce(slog.LevelWarn, msg.(string))
}

// LogOnceDebug logs a message with level debug only once.
func LogOnceDebug(msg any) {
	logOnce(slog.LevelDebug, msg.(string))
}

// LogOnChangeInfo logs a message with level info only when the message changes.
func LogOnChangeInfo(msg any) {
	logOnChange(slog.LevelInfo, msg.(string))
}

// LogOnChangeError logs a message with level error only when the message changes.
func LogOnChangeError(msg any) {
	logOnChange(slog.LevelError, msg.(string))
}

// LogOnChangeWarn logs a message with level warn only when the message changes.
func LogOnChangeWarn(msg any) {
	logOnChange(slog.LevelWarn, msg.(string))
}

// LogOnChangeDebug logs a message with level debug only when the message changes.
func LogOnChangeDebug(msg any) {
	logOnChange(slog.LevelDebug, msg.(string))
}

// Infof logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Infof(msg string, args ...any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelInfo, fmt.Sprintf(msg, args...), logcore.GetRecordInformation()))
}

// Errorf logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Errorf(msg string, args ...any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelError, fmt.Sprintf(msg, args...), logcore.GetRecordInformation()))
}

// Warnf logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warnf(msg string, args ...any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelWarn, fmt.Sprintf(msg, args...), logcore.GetRecordInformation()))
}

// Debugf logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debugf(msg string, args ...any) {
	logcore.Logger().SendReport(logcore.NewIonReport(slog.LevelDebug, fmt.Sprintf(msg, args...), logcore.GetRecordInformation()))
}

// LogOnceInfof logs a message with level info only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceInfof(msg string, args ...any) {
	logOnce(slog.LevelInfo, fmt.Sprintf(msg, args...))
}

// LogOnceErrorf logs a message with level error only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceErrorf(msg string, args ...any) {
	logOnce(slog.LevelError, fmt.Sprintf(msg, args...))
}

// LogOnceWarnf logs a message with level warn only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceWarnf(msg string, args ...any) {
	logOnce(slog.LevelWarn, fmt.Sprintf(msg, args...))
}

// LogOnceDebugf logs a message with level debug only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceDebugf(msg string, args ...any) {
	logOnce(slog.LevelDebug, fmt.Sprintf(msg, args...))
}

// LogOnChangeInfof logs a message with level info only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeInfof(msg string, args ...any) {
	logOnChange(slog.LevelInfo, fmt.Sprintf(msg, args...))
}

// LogOnChangeErrorf logs a message with level error only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeErrorf(msg string, args ...any) {
	logOnChange(slog.LevelError, fmt.Sprintf(msg, args...))
}

// LogOnChangeWarnf logs a message with level warn only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeWarnf(msg string, args ...any) {
	logOnChange(slog.LevelWarn, fmt.Sprintf(msg, args...))
}

// LogOnChangeDebugf logs a message with level debug only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeDebugf(msg string, args ...any) {
	logOnChange(slog.LevelDebug, fmt.Sprintf(msg, args...))
}

// logOnce logs a message with level info only once. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only once.
// Avoid using it in a sintax like this: LogOnce("Logging..."); LogOnce("Logging...")
func logOnce(level slog.Level, recordMsg string) {
	callInfo := logcore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	proceed := usecases.LogOnce(
		logcore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		logcore.Logger().SendReport(logcore.NewIonReport(level, recordMsg, callInfo))
	}
}

// logOnChange logs a message with level info only when the message changes. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only when it changes.
// Avoid using it in a sintax like this: LogOnChange("Logging..."); LogOnChange("Logging...")
func logOnChange(level slog.Level, recordMsg string) {
	callInfo := logcore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	proceed := usecases.LogOnChange(
		logcore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		logcore.Logger().SendReport(logcore.NewIonReport(level, recordMsg, callInfo))
	}
}
