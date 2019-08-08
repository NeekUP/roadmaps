package tests

import "fmt"

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
