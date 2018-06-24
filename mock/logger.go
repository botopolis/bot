package mock

// Level is the log level
type Level int

const (
	// DebugLevel is the debug level
	DebugLevel Level = iota
	// InfoLevel is the info level
	InfoLevel
	// WarnLevel is the warn level
	WarnLevel
	// ErrorLevel is the error level
	ErrorLevel
	// FatalLevel is the fatal level
	FatalLevel
	// PanicLevel is the panic level
	PanicLevel
)

// NewLogger creates a mock logger
func NewLogger() *Logger {
	return &Logger{
		WriteFunc:  func(Level, ...interface{}) {},
		WritefFunc: func(Level, string, ...interface{}) {},
	}
}

// Logger is a mock logger
type Logger struct {
	// WriteFunc intercepts non-formatted logs
	WriteFunc func(Level, ...interface{})
	// WritefFunc intercepts formatted logs
	WritefFunc func(Level, string, ...interface{})
}

// Debug logs at DebugLevel
func (l *Logger) Debug(v ...interface{}) { l.WriteFunc(DebugLevel, v...) }

// Debugf logs at DebugLevel with formatting
func (l *Logger) Debugf(fmt string, v ...interface{}) { l.WritefFunc(DebugLevel, fmt, v...) }

// Info logs at InfoLevel
func (l *Logger) Info(v ...interface{}) { l.WriteFunc(InfoLevel, v...) }

// Infof logs at InfoLevel with formatting
func (l *Logger) Infof(fmt string, v ...interface{}) { l.WritefFunc(InfoLevel, fmt, v...) }

// Warn logs at WarnLevel
func (l *Logger) Warn(v ...interface{}) { l.WriteFunc(WarnLevel, v...) }

// Warnf logs at WarnLevel with formatting
func (l *Logger) Warnf(fmt string, v ...interface{}) { l.WritefFunc(WarnLevel, fmt, v...) }

// Error logs at ErrorLevel
func (l *Logger) Error(v ...interface{}) { l.WriteFunc(ErrorLevel, v...) }

// Errorf logs at ErrorLevel with formatting
func (l *Logger) Errorf(fmt string, v ...interface{}) { l.WritefFunc(ErrorLevel, fmt, v...) }

// Fatal logs at FatalLevel
func (l *Logger) Fatal(v ...interface{}) { l.WriteFunc(FatalLevel, v...) }

// Fatalf logs at FatalLevel with formatting
func (l *Logger) Fatalf(fmt string, v ...interface{}) { l.WritefFunc(FatalLevel, fmt, v...) }

// Panic logs at PanicLevel
func (l *Logger) Panic(v ...interface{}) { l.WriteFunc(PanicLevel, v...) }

// Panicf logs at PanicLevel with formatting
func (l *Logger) Panicf(fmt string, v ...interface{}) { l.WritefFunc(PanicLevel, fmt, v...) }
