package core

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
