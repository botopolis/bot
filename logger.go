package bot

import (
	"log"
	"os"
)

// Logger is the interface Robot exposes
type Logger interface {
	Debug(values ...interface{})
	Debugf(format string, values ...interface{})

	Info(values ...interface{})
	Infof(format string, values ...interface{})

	Warn(values ...interface{})
	Warnf(format string, values ...interface{})

	Error(values ...interface{})
	Errorf(format string, values ...interface{})

	Fatal(values ...interface{})
	Fatalf(format string, values ...interface{})

	Panic(values ...interface{})
	Panicf(format string, values ...interface{})
}

// defaultLogger is a very basic logger that the bot loads by default
var defaultLogger = basicLogger{log.New(os.Stdout, "", 0)}

type basicLogger struct{ *log.Logger }

func (l basicLogger) Debug(v ...interface{})            {}
func (l basicLogger) Debugf(f string, v ...interface{}) {}
func (l basicLogger) Info(v ...interface{})             { l.Logger.Println(v...) }
func (l basicLogger) Infof(f string, v ...interface{})  { l.Logger.Printf(f, v...) }
func (l basicLogger) Warn(v ...interface{})             { l.Logger.Println(v...) }
func (l basicLogger) Warnf(f string, v ...interface{})  { l.Logger.Printf(f, v...) }
func (l basicLogger) Error(v ...interface{})            { l.Logger.Println(v...) }
func (l basicLogger) Errorf(f string, v ...interface{}) { l.Logger.Printf(f, v...) }
func (l basicLogger) Panic(v ...interface{})            { l.Logger.Panicln(v...) }
func (l basicLogger) Panicf(f string, v ...interface{}) { l.Logger.Panicf(f, v...) }
