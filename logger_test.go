package bot_test

import (
	"fmt"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
)

type logStuff struct{}

func (logStuff) Load(r *bot.Robot) {
	r.Logger.Error("One")
	r.Logger.Panic("Two")
	r.Logger.Fatal("Three")
}

func ExampleLogger() {
	// Ignore this - just example setup
	chat := mock.NewChat()
	chat.MessageChan = make(chan bot.Message)
	go func() { close(chat.MessageChan) }()

	logger := mock.NewLogger()
	logger.WriteFunc = func(level mock.Level, v ...interface{}) {
		switch level {
		case mock.ErrorLevel:
			fmt.Println("Error log")
		case mock.PanicLevel:
			fmt.Println("Panic log")
		case mock.FatalLevel:
			fmt.Println("Fatal log")

		}
	}
	b := bot.New(chat, logStuff{})
	b.Logger = logger
	b.Run()
	// Output:
	// Error log
	// Panic log
	// Fatal log
}
