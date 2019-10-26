package logger

// Log is a package level variable, every program should access logging function through "Log"
var log Logger

//Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

// Logger represent common interface for logging function
type Logger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})

	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})

	Infof(format string, args ...interface{})
	Info(args ...interface{})

	Warnf(format string, args ...interface{})
	Warn(args ...interface{})

	Debugf(format string, args ...interface{})
	Debug(args ...interface{})

	Panicf(format string, args ...interface{})
	Panic(args ...interface{})

	WithFields(keyValues Fields) Logger
}

//InitLogger returns an instance of logger
func InitLogger() error {
	logger, err := NewZapLogger()
	if err != nil {
		return err
	}
	log = logger
	return nil
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func Panic(args ...interface{}) {
	log.Panic(args...)
}

func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

var _ Logger = (*zapLogger)(nil)