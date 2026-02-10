package logger

import "context"

type contextKey string

const loggerKey contextKey = "logger"

func ToContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if loggger, ok := ctx.Value(loggerKey).(Logger); ok {
		return loggger
	}

	return NewNoOp()
}
