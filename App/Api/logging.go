package lostpets

// Logger logger interface
type Logger interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	UnwrapError(err error)
}

//counterfeiter:generate -o mock/mockStructuredLogger.go --fake-name StructuredLogger . StructuredLogger
type StructuredLogger interface {
	WithFields(fields map[string]interface{}) Logger
	Logger
}
