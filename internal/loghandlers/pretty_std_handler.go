package loghandlers

import (
	"context"
	"io"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyStdLogHandler struct {
	Writer io.Writer
	Level  slog.Level
}

func (h *PrettyStdLogHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := r.Time.Format("15:04:05.00")
	var levelColor *color.Color
	var levelStr string

	switch r.Level {
	case slog.LevelDebug:
		levelColor = color.New(color.FgHiBlack)
		levelStr = "D"
	case slog.LevelInfo:
		levelColor = color.New(color.FgGreen)
		levelStr = "I"
	case slog.LevelWarn:
		levelColor = color.New(color.FgYellow)
		levelStr = "W"
	case slog.LevelError:
		levelColor = color.New(color.FgRed)
		levelStr = "E"
	default:
		levelColor = color.New(color.FgCyan)
		levelStr = r.Level.String()
	}

	_, err := io.WriteString(h.Writer, "[")
	if err != nil {
		return err
	}

	_, err = levelColor.Fprint(h.Writer, levelStr)
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

func (h *PrettyStdLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyStdLogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *PrettyStdLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level
}
