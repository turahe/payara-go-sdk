package payara

// NopLogger is a no-op logger that implements Logger. Use when no logger is provided.
type NopLogger struct{}

func (NopLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (NopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (NopLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (NopLogger) Error(msg string, keysAndValues ...interface{})  {}

// Ensure NopLogger implements Logger
var _ Logger = (*NopLogger)(nil)
