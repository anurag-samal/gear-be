package logger

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
)

const logDir = "logs"

type Options struct {
	Level     slog.Level
	AddSource bool
	JSON      bool
}

func DefaultOptions(env string) Options {
	return Options{
		Level:     envLogLevel(env),
		AddSource: env == "development",
		JSON:      env == "production",
	}
}

func Init(opts Options) error {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	errorFile, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	combinedFile, err := os.OpenFile(filepath.Join(logDir, "combined.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		errorFile.Close()
		return err
	}

	consoleHandler := newConsoleHandler(opts)

	handler := &multiHandler{
		handlers: []slog.Handler{
			consoleHandler,
			slog.NewJSONHandler(combinedFile, &slog.HandlerOptions{Level: opts.Level}),
			slog.NewJSONHandler(errorFile, &slog.HandlerOptions{Level: slog.LevelError}),
		},
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

func InitDefault(env string) {
	opts := DefaultOptions(env)
	if err := Init(opts); err != nil {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})))
	}
}

func envLogLevel(env string) slog.Level {
	switch env {
	case "production":
		return slog.LevelInfo
	default:
		return slog.LevelDebug
	}
}

func newConsoleHandler(opts Options) slog.Handler {
	ho := &slog.HandlerOptions{
		Level:     opts.Level,
		AddSource: opts.AddSource,
	}

	if opts.JSON {
		return slog.NewJSONHandler(os.Stdout, ho)
	}
	return slog.NewTextHandler(os.Stdout, ho)
}

type multiHandler struct {
	handlers []slog.Handler
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var firstErr error
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: handlers}
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: handlers}
}

var _ slog.Handler = (*multiHandler)(nil)
