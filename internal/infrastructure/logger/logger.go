package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	// Context-aware logging
	WithContext(ctx context.Context) Logger
	With(keysAndValues ...interface{}) Logger
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func New(level string) Logger {
	config := zap.NewProductionConfig()

	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config.Encoding = "json"
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.StacktraceKey = "stacktrace"

	baseLogger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("failed to initialiaze logger %v", err))
	}

	return &zapLogger{
		logger: baseLogger.Sugar(),
	}
}

func NewDevelopment() Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	baseLogger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger %v", err))
	}

	return &zapLogger{
		logger: baseLogger.Sugar(),
	}
}

// Debug implements [Logger].
func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.logger.Debugw(msg, keysAndValues...)

}

// Error implements [Logger].
func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

// Fatal implements [Logger].
func (z *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	z.logger.Fatalw(msg, keysAndValues...)
}

// Info implements [Logger].
func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Infow(msg, keysAndValues...)
}

// Warn implements [Logger].
func (z *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	z.logger.Warnw(msg, keysAndValues...)
}

// With implements [Logger].
func (z *zapLogger) With(keysAndValues ...interface{}) Logger {
	return &zapLogger{
		logger: z.logger.With(keysAndValues...),
	}
}

// WithContext implements [Logger].
func (z *zapLogger) WithContext(ctx context.Context) Logger {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return &zapLogger{
			logger: z.logger.With("request_id", requestID),
		}
	}

	return z
}
