package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger to provide structured logging
type Logger struct {
	zap *zap.Logger
}

// Fields is a map of key-value pairs for structured logging
type Fields map[string]interface{}

// Global logger instance
var globalLogger *Logger

// Initialize creates and configures the global logger
func Initialize(environment string, level string) error {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(zapLevel)

	// Build logger
	zapLogger, err := config.Build(
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	globalLogger = &Logger{zap: zapLogger}
	return nil
}

// Get returns the global logger instance
func Get() *Logger {
	if globalLogger == nil {
		// Fallback to a basic logger if not initialized
		zapLogger, _ := zap.NewProduction()
		globalLogger = &Logger{zap: zapLogger}
	}
	return globalLogger
}

// WithContext returns a logger with context values
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract correlation ID from context if present
	if correlationID, ok := ctx.Value("correlation_id").(string); ok {
		return &Logger{
			zap: l.zap.With(zap.String("correlation_id", correlationID)),
		}
	}
	return l
}

// WithFields returns a logger with additional fields
func (l *Logger) WithFields(fields Fields) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{
		zap: l.zap.With(zapFields...),
	}
}

// WithField returns a logger with an additional field
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		zap: l.zap.With(zap.Any(key, value)),
	}
}

// WithError returns a logger with an error field
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		zap: l.zap.With(zap.Error(err)),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string) {
	l.zap.Debug(msg)
}

// Debugf logs a formatted debug message
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zap.Sugar().Debugf(format, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string) {
	l.zap.Info(msg)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.zap.Sugar().Infof(format, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.zap.Warn(msg)
}

// Warnf logs a formatted warning message
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zap.Sugar().Warnf(format, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.zap.Error(msg)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zap.Sugar().Errorf(format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string) {
	l.zap.Fatal(msg)
}

// Fatalf logs a formatted fatal message and exits
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zap.Sugar().Fatalf(format, args...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// Package-level convenience functions that use the global logger

// Debug logs a debug message using the global logger
func Debug(msg string) {
	Get().Debug(msg)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(format string, args ...interface{}) {
	Get().Debugf(format, args...)
}

// Info logs an info message using the global logger
func Info(msg string) {
	Get().Info(msg)
}

// Infof logs a formatted info message using the global logger
func Infof(format string, args ...interface{}) {
	Get().Infof(format, args...)
}

// Warn logs a warning message using the global logger
func Warn(msg string) {
	Get().Warn(msg)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(format string, args ...interface{}) {
	Get().Warnf(format, args...)
}

// Error logs an error message using the global logger
func Error(msg string) {
	Get().Error(msg)
}

// Errorf logs a formatted error message using the global logger
func Errorf(format string, args ...interface{}) {
	Get().Errorf(format, args...)
}

// Fatal logs a fatal message and exits using the global logger
func Fatal(msg string) {
	Get().Fatal(msg)
}

// Fatalf logs a formatted fatal message and exits using the global logger
func Fatalf(format string, args ...interface{}) {
	Get().Fatalf(format, args...)
}

// WithFields returns a logger with additional fields using the global logger
func WithFields(fields Fields) *Logger {
	return Get().WithFields(fields)
}

// WithField returns a logger with an additional field using the global logger
func WithField(key string, value interface{}) *Logger {
	return Get().WithField(key, value)
}

// WithError returns a logger with an error field using the global logger
func WithError(err error) *Logger {
	return Get().WithError(err)
}

// WithContext returns a logger with context values using the global logger
func WithContext(ctx context.Context) *Logger {
	return Get().WithContext(ctx)
}

// Sync flushes any buffered log entries from the global logger
func Sync() error {
	return Get().Sync()
}

// NewNopLogger creates a no-op logger for testing
func NewNopLogger() *Logger {
	return &Logger{zap: zap.NewNop()}
}

// SetOutput redirects log output to the specified writer (mainly for testing)
func SetOutput(writer zapcore.WriteSyncer) error {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	zapLogger := zap.New(core)
	globalLogger = &Logger{zap: zapLogger}
	return nil
}

// GetZapLogger returns the underlying zap logger (for advanced use cases)
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zap
}

// ExitFunc is used for testing Fatal functions
var ExitFunc = os.Exit
