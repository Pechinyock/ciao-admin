package loghandlers

import (
	"context"
	"fmt"
	"log/slog"
)

type MultiOutputHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiOutputHandler {
	if len(handlers) == 0 {
		return nil
	}
	result := &MultiOutputHandler{
		handlers: handlers,
	}

	return result
}

func (h *MultiOutputHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for i, handler := range h.handlers {
		if handler.Enabled(ctx, r.Level) {
			if err := handler.Handle(ctx, r); err != nil {
				errs = append(errs, fmt.Errorf("handler %d: %w", i, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("multiHandler errors: %v", errs)
	}

	return nil
}

func (h *MultiOutputHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *MultiOutputHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *MultiOutputHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}
