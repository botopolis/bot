package gobot_test

import (
	"fmt"
	"net/http"

	"github.com/berfarah/gobot"
)

type ExamplePlugin struct {
	Username string
}

func (p *ExamplePlugin) Load(r *gobot.Robot) { p.Username = "beardroid" }

func Example() {
	robot := gobot.New(NewChat())

	// Install plugins
	robot.Install(
		&ExamplePlugin{},
	)

	// Respond to any message
	robot.Hear(gobot.Regexp("welcome (\\w*)"), func(r gobot.Responder) error {
		return r.Send(gobot.Message{
			Text: "All hail " + r.Match[1] + ", our new overlord",
		})
	})

	// Respond to messages that are either DMed to the bot or start with the bot's name
	robot.Respond(gobot.Regexp("topic (.*)"), func(r gobot.Responder) error {
		return r.Topic(r.Match[1])
	})

	// Track when topics are updated
	robot.Topic(func(r gobot.Responder) error {
		if r.Room == "announcements" {
			var announcements []string
			r.Brain.Get("announcements", &announcements)
			announcements = append(announcements, r.Match[1])
			r.Brain.Set("announcements", &announcements)
		}
		return nil
	})

	// Track the comings and goings of people
	robot.Enter(func(r gobot.Responder) error { return nil })
	robot.Leave(func(r gobot.Responder) error { return nil })

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
		// "beardroid"
	}
}
