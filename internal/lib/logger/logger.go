package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	formatText = "text"
	formatJSON = "json"
)

func SetupLogger() *slog.Logger {
	format := viper.GetString("LOG_FORMAT")
	opts := &slog.HandlerOptions{
		Level: ParseLevel(viper.GetString("LOG_LEVEL")),
	}
	var log slog.Handler

	switch format {
	case formatText:
		log = slog.NewTextHandler(os.Stdout, opts)
	case formatJSON:
		log = slog.NewJSONHandler(os.Stdout, opts)
	default:
		log = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(log)
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
