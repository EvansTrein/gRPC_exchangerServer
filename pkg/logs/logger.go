package logs

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"time"
)

type CustomHandler struct {
	handler slog.Handler
	output  io.Writer
}

func NewCustomHandler(output io.Writer, opts *slog.HandlerOptions) *CustomHandler {
	return &CustomHandler{
		handler: slog.NewTextHandler(output, opts),
		output:  output,
	}
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	h.output.Write([]byte("\n"))

	h.output.Write([]byte("Apps log:\n"))

	h.output.Write([]byte(r.Time.Format(time.Stamp) + "\n"))

	level := r.Level.String()
	switch r.Level {
	case slog.LevelInfo:
		level = "\033[32m" + level + "\033[0m" // green color
	case slog.LevelError:
		level = "\033[31m" + level + "\033[0m" // red color
	case slog.LevelDebug:
		level = "\033[34m" + level + "\033[0m" // blue color
	case slog.LevelWarn:
		level = "\033[33m" + level + "\033[0m" // yellow color
	}
	h.output.Write([]byte("level--> " + level + "\n"))

	h.output.Write([]byte("msg--> " + r.Message + "\n"))

	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		source := "file--> " + f.File + "\ncode_line--> " + strconv.Itoa(f.Line) + "\n"
		h.output.Write([]byte(source))
	}

	r.Attrs(func(attr slog.Attr) bool {
		h.output.Write([]byte(attr.Key + "--> " + attr.Value.String() + "\n"))
		return true
	})

	h.output.Write([]byte("\n"))
	return nil
}

func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithAttrs(attrs),
		output:  h.output,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		handler: h.handler.WithGroup(name),
		output:  h.output,
	}
}

func InitLog(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(NewCustomHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, AddSource: true}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
