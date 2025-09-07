package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type LogFormat string

const (
	FormatText LogFormat = "text"
	FormatMin  LogFormat = "min"
	FormatJSON LogFormat = "json"
)

// InitFromEnv configures the global slog logger based on env and returns the LevelVar.
// Supported LOG_LEVEL: debug | info | warn | error (case-insensitive).
// Uses JSON output and includes source (file:line).
func InitFromEnv(logLevel string, logFormat string) *slog.LevelVar {

	levelVar := new(slog.LevelVar)
	levelVar.Set(slog.LevelInfo) // by default

	levelVar.Set(parseLevel(logLevel))

	format, ok := parseFormat(logFormat)
	if !ok {
		format = FormatJSON
	}

	var handler slog.Handler
	switch format {
	case FormatText:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     levelVar,
			AddSource: true,
		})
	case FormatMin:
		handler = newMinimalHandler(os.Stdout, levelVar)
	default: // FormatJSON
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     levelVar,
			AddSource: true,
		})
	}

	slog.SetDefault(slog.New(handler))

	slog.Info("Logger initialized", "level", logLevel)

	return levelVar
}

func parseLevel(logLevel string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(logLevel)) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	}
	return slog.LevelInfo
}

func parseFormat(s string) (LogFormat, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text":
		return FormatText, true
	case "min":
		return FormatMin, true
	case "", "json":
		return FormatJSON, true
	default:
		return FormatJSON, false
	}
}

// ---------------- minimal handler ----------------

// minimalHandler prints only the message (one line), ignoring attrs/time/source.
type minimalHandler struct {
	w   io.Writer
	lvl *slog.LevelVar
	mu  sync.Mutex
}

func newMinimalHandler(w io.Writer, lvl *slog.LevelVar) *minimalHandler {
	return &minimalHandler{w: w, lvl: lvl}
}

func (h *minimalHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.lvl.Level()
}

func (h *minimalHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.w, r.Message+"\n")
	return err
}

func (h *minimalHandler) WithAttrs(_ []slog.Attr) slog.Handler { return h }
func (h *minimalHandler) WithGroup(_ string) slog.Handler      { return h }
