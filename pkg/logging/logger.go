package logging

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface for structured logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger
}

// Field represents a log field
type Field = zapcore.Field

// Common field constructors
var (
	String   = zap.String
	Int      = zap.Int
	Int64    = zap.Int64
	Float64  = zap.Float64
	Bool     = zap.Bool
	Time     = zap.Time
	Duration = zap.Duration
	Error    = zap.Error
	Any      = zap.Any
	Strings  = zap.Strings
	Ints     = zap.Ints
)

// zapLogger wraps zap.Logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
}

// Config contains logger configuration
type Config struct {
	Level       string // debug, info, warn, error, fatal
	Format      string // json, console
	Output      string // stdout, stderr, or file path
	Development bool
	AddCaller   bool
}

// NewLogger creates a new structured logger
func NewLogger(cfg Config) (Logger, error) {
	// Parse level
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Configure encoder
	var encoderConfig zapcore.EncoderConfig
	if cfg.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	// Customize encoder config
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create encoder
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Configure output
	var writer io.Writer
	switch cfg.Output {
	case "stdout":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// File output
		file, err := os.OpenFile(cfg.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		writer = file
	}

	// Create core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)

	// Create logger options
	opts := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if cfg.AddCaller {
		opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	// Create logger
	logger := zap.New(core, opts...)

	return &zapLogger{logger: logger}, nil
}

// NewDevelopmentLogger creates a logger for development
func NewDevelopmentLogger() (Logger, error) {
	return NewLogger(Config{
		Level:       "debug",
		Format:      "console",
		Output:      "stdout",
		Development: true,
		AddCaller:   true,
	})
}

// NewProductionLogger creates a logger for production
func NewProductionLogger() (Logger, error) {
	return NewLogger(Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		Development: false,
		AddCaller:   false,
	})
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
	}
}

// WithContext creates a logger with trace context
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return l
	}

	spanContext := span.SpanContext()
	return l.With(
		String("trace_id", spanContext.TraceID().String()),
		String("span_id", spanContext.SpanID().String()),
	)
}

// ContextKey is the key type for logger in context
type contextKey string

const loggerKey contextKey = "logger"

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext retrieves logger from context
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	// Return no-op logger if not found
	return &noopLogger{}
}

// noopLogger implements Logger but does nothing
type noopLogger struct{}

func (n *noopLogger) Debug(msg string, fields ...Field) {}
func (n *noopLogger) Info(msg string, fields ...Field)  {}
func (n *noopLogger) Warn(msg string, fields ...Field)  {}
func (n *noopLogger) Error(msg string, fields ...Field) {}
func (n *noopLogger) Fatal(msg string, fields ...Field) {}
func (n *noopLogger) With(fields ...Field) Logger       { return n }
func (n *noopLogger) WithContext(ctx context.Context) Logger { return n }

// Helper functions for common log patterns

// LogRequest logs an HTTP request
func LogRequest(logger Logger, method, path string, statusCode int, duration time.Duration, fields ...Field) {
	allFields := append(fields,
		String("method", method),
		String("path", path),
		Int("status", statusCode),
		Duration("duration", duration),
	)

	level := "info"
	if statusCode >= 500 {
		level = "error"
	} else if statusCode >= 400 {
		level = "warn"
	}

	switch level {
	case "error":
		logger.Error("HTTP request", allFields...)
	case "warn":
		logger.Warn("HTTP request", allFields...)
	default:
		logger.Info("HTTP request", allFields...)
	}
}

// LogError logs an error with context
func LogError(logger Logger, ctx context.Context, msg string, err error, fields ...Field) {
	allFields := append(fields, Error(err))
	logger.WithContext(ctx).Error(msg, allFields...)
}