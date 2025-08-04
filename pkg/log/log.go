package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
)

var defaultLogger *slog.Logger

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

type Config struct {
	Level  LogLevel
	Output io.Writer
}

func DefaultConfig() Config {
	return Config{
		Level:  getLogLevelFromEnv(),
		Output: os.Stdout,
	}
}

func Initialize(config Config) {
	var level slog.Level
	switch config.Level {
	case LevelDebug:
		level = slog.LevelDebug
	case LevelInfo:
		level = slog.LevelInfo
	case LevelWarn:
		level = slog.LevelWarn
	case LevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: config.Level == LevelDebug,
	}

	var handler = slog.NewTextHandler(config.Output, opts)
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func InitializeDefault() {
	Initialize(DefaultConfig())
}

func Debug(msg string, args ...any) {
	defaultLogger.Debug(msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.DebugContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	defaultLogger.Info(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.InfoContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	defaultLogger.Warn(msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.WarnContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
}

func With(args ...any) *slog.Logger {
	return defaultLogger.With(args...)
}

func WithGroup(name string) *slog.Logger {
	return defaultLogger.WithGroup(name)
}

func GetLogger() *slog.Logger {
	return defaultLogger
}

func Fatal(msg string, args ...any) {
	defaultLogger.Error(msg, args...)
	os.Exit(1)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	defaultLogger.ErrorContext(ctx, msg, args...)
	os.Exit(1)
}

func getLogLevelFromEnv() LogLevel {
	level := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn", "warning":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
