package logging

import (
	"log/slog"
	"path"
	"runtime"
)

// WithPackage returns a logger bound to the caller's package name as "component".
func WithPackage() *slog.Logger {
	return slog.With("component", callerPackage(1))
}

// Package returns the caller's package name (plain string).
func Package() string {
	return callerPackage(1)
}

// --- internals ---

// callerPackage extracts package name from the caller's file path.
func callerPackage(skip int) string {
	_, file, _, ok := runtime.Caller(skip + 1)
	if !ok || file == "" {
		return "unknown"
	}

	return path.Base(path.Dir(file))
}
