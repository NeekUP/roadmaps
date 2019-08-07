package infrastructure

import "fmt"

type AppLogger interface {
	// Debug uses fmt.Sprint to construct and log a message.
	Debug(args ...interface{})

	// Info uses fmt.Sprint to construct and log a message.
	Info(args ...interface{})

	// Warn uses fmt.Sprint to construct and log a message.
	Warn(args ...interface{})

	// Error uses fmt.Sprint to construct and log a message.
	Error(args ...interface{})

	// DPanic uses fmt.Sprint to construct and log a message. In development, the
	// logger then panics. (See DPanicLevel for details.)
	DPanic(args ...interface{})

	// Panic uses fmt.Sprint to construct and log a message, then panics.
	Panic(args ...interface{})

	// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
	Fatal(args ...interface{})

	// Debugf uses fmt.Sprintf to log a templated message.
	Debugf(template string, args ...interface{})

	// Infof uses fmt.Sprintf to log a templated message.
	Infof(template string, args ...interface{})

	// Warnf uses fmt.Sprintf to log a templated message.
	Warnf(template string, args ...interface{})

	// Errorf uses fmt.Sprintf to log a templated message.
	Errorf(template string, args ...interface{})

	// DPanicf uses fmt.Sprintf to log a templated message. In development, the
	// logger then panics. (See DPanicLevel for details.)
	DPanicf(template string, args ...interface{})

	// Panicf uses fmt.Sprintf to log a templated message, then panics.
	Panicf(template string, args ...interface{})

	// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
	Fatalf(template string, args ...interface{})

	// When debug-level logging is disabled, this is much faster than
	//  s.With(keysAndValues).Debug(msg)
	Debugw(msg string, keysAndValues ...interface{})

	// Infow logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Infow(msg string, keysAndValues ...interface{})

	// Warnw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Warnw(msg string, keysAndValues ...interface{})

	// Errorw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Errorw(msg string, keysAndValues ...interface{})

	// DPanicw logs a message with some additional context. In development, the
	// logger then panics. (See DPanicLevel for details.) The variadic key-value
	// pairs are treated as they are in With.
	DPanicw(msg string, keysAndValues ...interface{})

	// Panicw logs a message with some additional context, then panics. The
	// variadic key-value pairs are treated as they are in With.
	Panicw(msg string, keysAndValues ...interface{})
}

type AppLoggerForTests struct{}

func (AppLoggerForTests) Debug(args ...interface{}) {

}

func (AppLoggerForTests) Info(args ...interface{}) {
}

func (AppLoggerForTests) Warn(args ...interface{}) {
}

func (AppLoggerForTests) Error(args ...interface{}) {
}

func (AppLoggerForTests) DPanic(args ...interface{}) {
}

func (AppLoggerForTests) Panic(args ...interface{}) {
}

func (AppLoggerForTests) Fatal(args ...interface{}) {
}

func (AppLoggerForTests) Debugf(template string, args ...interface{}) {
}

func (AppLoggerForTests) Infof(template string, args ...interface{}) {
}

func (AppLoggerForTests) Warnf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (AppLoggerForTests) Errorf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (AppLoggerForTests) DPanicf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (AppLoggerForTests) Panicf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (AppLoggerForTests) Fatalf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (AppLoggerForTests) Debugw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (AppLoggerForTests) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (AppLoggerForTests) Warnw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (AppLoggerForTests) Errorw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (AppLoggerForTests) DPanicw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (AppLoggerForTests) Panicw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}
