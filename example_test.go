package bot_test

import (
	"fmt"
	"net/http"

	"github.com/botopolis/bot"
	"github.com/botopolis/bot/mock"
)

type ExamplePlugin struct {
	Username string
}

func (p *ExamplePlugin) Load(r *bot.Robot) { p.Username = "beardroid" }

func Example() {
	// Ignore this - just example setup
	chat := mock.NewChat()
	chat.MessageChan = make(chan bot.Message)
	go func() { close(chat.MessageChan) }()

	// Install adapter and plugins
	robot := bot.New(
		chat,
		&ExamplePlugin{},
	)

	// Respond to any message
	robot.Hear(bot.Regexp("welcome (\\w*)"), func(r bot.Responder) error {
		return r.Send(bot.Message{
			Text: "All hail " + r.Match[1] + ", our new overlord",
		})
	})

	// Respond to messages that are either DMed to the bot or start with the bot's name
	robot.Respond(bot.Regexp("topic (.*)"), func(r bot.Responder) error {
		return r.Topic(r.Match[1])
	})

	// Track when topics are updated
	robot.Topic(func(r bot.Responder) error {
		if r.Room == "announcements" {
			var announcements []string
			r.Brain.Get("announcements", &announcements)
			announcements = append(announcements, r.Match[1])
			r.Brain.Set("announcements", &announcements)
		}
		return nil
	})

	// Track the comings and goings of people
	robot.Enter(func(r bot.Responder) error { return nil })
	robot.Leave(func(r bot.Responder) error { return nil })

	// Create server endpoints
	robot.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})

	// Start the server and listening
	robot.Run()

	// Access plugins
	var plugin ExamplePlugin
	if ok := robot.Plugin(&plugin); ok {
		fmt.Println(plugin.Username)
		// Output: beardroid
	}
}
