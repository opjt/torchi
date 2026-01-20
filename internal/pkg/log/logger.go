package log

import (
	"log/slog"
	"os"
	"torchi/internal/pkg/config"
)

// Logger Struct
type Logger struct {
	*slog.Logger
}

// parseLevel : string → slog.Level
func parseLevel(levelStr string) slog.Level {
	switch levelStr {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewServiceLogger : for stdout , dev/prod
func NewLogger(env config.Env) *Logger {
	// stdout handler
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: parseLevel(env.Log.Level),
	})

	if env.Stage == config.StageDev {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     parseLevel(env.Log.Level),
			AddSource: true, // caller 정보 포함
		})
	}

	rawLogger := slog.New(handler)
	return &Logger{Logger: rawLogger}
}

func (l *Logger) Fatal(err error) {
	l.Error(err.Error())
	os.Exit(1)
}
