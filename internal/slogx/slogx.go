package slogx

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func Fatal(msg string, err error) {
	slog.Error(msg, "error", err.Error())
	os.Exit(1)
}

// NewHandler returns slog.Handler corresponding to given format (text or json),
// slog.TextHandler will be returned on invalid format.
func NewHandler(w io.Writer, format, level string, addSource bool) slog.Handler {
	opts := &slog.HandlerOptions{
		AddSource: addSource,
		Level:     parseLevel(level),
	}

	if strings.ToLower(format) == "json" {
		return slog.NewJSONHandler(w, opts)
	}

	return slog.NewTextHandler(w, opts)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error", "err":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
