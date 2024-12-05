package ionlog

import (
	"io"
	"log/slog"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
)

type customAttrs func(i ioncore.IIonLogger)

type periodicRotation int

const (
	// Daily rotation
	Daily periodicRotation = iota
	// Weekly rotation
	Weekly
	// Monthly rotation
	Monthly
)

// SetLogAttributes sets the log SetLogAttributes
// fns is a variadic parameter that accepts customAttrs
func SetLogAttributes(fns ...customAttrs) {
	for _, fn := range fns {
		fn(ioncore.Logger())
	}
}

// WithTargets sets the write targets for the logger, every log will be written to these targets.
func WithTargets(w ...io.Writer) customAttrs {
	return func(i ioncore.IIonLogger) {
		i.SetTargets(w...)
	}
}

// WithStaicFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaicFields(attrs map[string]string) customAttrs {
	return func(i ioncore.IIonLogger) {
		index := 0
		slogAttrs := make([]slog.Attr, len(attrs))
		for k, v := range attrs {
			slogAttrs[index] = slog.String(k, v)
			index++
		}
		handler := i.CreateDefaultLogHandler().WithAttrs(slogAttrs)
		i.SetLogEngine(slog.New(handler))
	}
}

func WithLogFileRotation() customAttrs {
	return func(i ioncore.IIonLogger) {

	}
}
