package logger

import "context"

type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (n *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {}
func (n *NoOpLogger) Fatal(msg string, keysAndValues ...interface{}) {}
func (n *NoOpLogger) Info(msg string, keysAndValues ...interface{})  {}
func (n *NoOpLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (n *NoOpLogger) With(keysAndValues ...interface{}) Logger       { return n }
func (n *NoOpLogger) WithContext(ctx context.Context) Logger         { return n }

func NewNoOp() Logger {
	return &NoOpLogger{}
}
