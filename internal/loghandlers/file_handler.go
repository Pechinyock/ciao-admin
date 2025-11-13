package loghandlers

import (
	"context"
	"io"
	"log/slog"
)

type FileLogHandler struct {
	Writer io.Writer
	Level  slog.Level
}

func (h *FileLogHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := r.Time.Format("15:04:05.00")
	var levelStr string

	switch r.Level {
	case slog.LevelDebug:
		levelStr = "D"
	case slog.LevelInfo:
		levelStr = "I"
	case slog.LevelWarn:
		levelStr = "W"
	case slog.LevelError:
		levelStr = "E"
	default:
		levelStr = r.Level.String()
	}

	_, err := io.WriteString(h.Writer, "[")
	if err != nil {
		return err
	}

	_, err = io.WriteString(h.Writer, levelStr)
	if err != nil {
		return err
	}
	_, err = io.WriteString(h.Writer, "] ["+timestamp+"]: "+r.Message)
	if err != nil {
		return err
	}

	if r.NumAttrs() > 0 {
		_, err = io.WriteString(h.Writer, " [")
		if err != nil {
			return err
		}
		first := true
		r.Attrs(func(attr slog.Attr) bool {
			if !first {
				_, err = io.WriteString(h.Writer, " | ")
				if err != nil {
					return false
				}
			}
			_, err = io.WriteString(h.Writer, attr.Key+": "+attr.Value.String())
			if err != nil {
				return false
			}
			first = false
			return true
		})
		_, err = io.WriteString(h.Writer, "]")
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(h.Writer, "\n")
	return err
}

func (h *FileLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *FileLogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *FileLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level
}
