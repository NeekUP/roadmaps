package tests

import "fmt"

type appLoggerForTests struct{}

func (appLoggerForTests) Debug(args ...interface{}) {

}

func (appLoggerForTests) Info(args ...interface{}) {
}

func (appLoggerForTests) Warn(args ...interface{}) {
}

func (appLoggerForTests) Error(args ...interface{}) {
}

func (appLoggerForTests) DPanic(args ...interface{}) {
}

func (appLoggerForTests) Panic(args ...interface{}) {
}

func (appLoggerForTests) Fatal(args ...interface{}) {
}

func (appLoggerForTests) Debugf(template string, args ...interface{}) {
}

func (appLoggerForTests) Infof(template string, args ...interface{}) {
}

func (appLoggerForTests) Warnf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (appLoggerForTests) Errorf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (appLoggerForTests) DPanicf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (appLoggerForTests) Panicf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (appLoggerForTests) Fatalf(template string, args ...interface{}) {
	fmt.Println(template, args)
}

func (appLoggerForTests) Debugw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (appLoggerForTests) Infow(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (appLoggerForTests) Warnw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (appLoggerForTests) Errorw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (appLoggerForTests) DPanicw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}

func (appLoggerForTests) Panicw(msg string, keysAndValues ...interface{}) {
	fmt.Println(msg, keysAndValues)
}
