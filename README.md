# botopolis

[![GoDoc](https://godoc.org/github.com/botopolis/bot?status.svg)](https://godoc.org/github.com/botopolis/bot) [![Build Status](https://circleci.com/gh/botopolis/bot.svg?style=svg)](https://circleci.com/gh/botopolis/bot) [![Test Coverage](https://api.codeclimate.com/v1/badges/b7acc61121363e7405a3/test_coverage)](https://codeclimate.com/github/botopolis/bot/test_coverage)

A hubot clone in Go! botopolis is extendable with plugins and works with different
chat services.

## Usage

See [example_test.go](./example_test.go) for usage details.

### Example program

Here's an example of a program you can write with `botopolis`. If you've used
[Hubot](https://github.com/hubotio/hubot) before, you'll see a familiar API.

```go
package main

import (
	"github.com/botopolis/bot"
	"github.com/botopolis/slack"
)

func main() {
	// Create a bot with a Slack adapter
	robot := bot.New(
		slack.New(os.Getenv("SLACK_TOKEN")),
	)

	// Create a Listener for messages to the bot with the word "hi"
	r.Respond(bot.Regexp("hi", func(r bot.Responder) error {
		// Respond to the sender of the message
		r.Reply("hello to you too!")
	})

	// Run the bot
	r.Run()
}
```

### Listeners

You have the following listeners at your disposal to interact with incoming
messages:

```go
// Listens for any message
robot.Hear(bot.Matcher, func(bot.Responder) error)
// Listens for messages addressed to the bot
robot.Respond(bot.Matcher, func(bot.Responder) error)
// Listens for topic changes
robot.Topic(func(bot.Responder) error)
// Listens for people entering a channel
robot.Enter(func(bot.Responder) error)
// Listens for people leaving a channel
robot.Leave(func(bot.Responder) error)
```

### Matchers

You'll notice that some listeners take a
[`bot.Matcher`](https://godoc.org/github.com/botopolis/bot#Matcher). Matchers
will run before the callback to determine whether the callback should be fired.
There are a few matchers provided by `botopolis` and you can write your own as
long as it matches the function signature.

```go
// Match the text of the message against a regular expression
bot.Regexp("^hi")
// Match a subset of the text
bot.Contains("hello")
// Match the user of a message
bot.User("jean")
// Match the room of a message
bot.Room("general")
```

### Callback and Responder

Callbacks are given a
[`bot.Responder`](https://godoc.org/github.com/botopolis/bot#Responder<Paste>).
The Responder holds a reference to the incoming `bot.Message` and exposes a few
methods which simplify interacting with the chat service.

```go
type Responder struct {
	Name string
	User string
	Room string
	Text string
	// abridged
}

// Send message
func (r Responder) Send(Message) error
// DM message
func (r Responder) Direct(string) error
// Reply to message
func (r Responder) Reply(string) error
// Change topic
func (r Responder) Topic(string) error
```

### Brain

[`bot.Brain`](https://godoc.org/github.com/botopolis/bot#Brain) is a store with
customizable backends. By default it will only save things in memory, but you
can add or create your own stores such as
[`botopolis/redis`](https://github.com/botopolis/redis) for data persistence.

```go
r := bot.New(mock.NewChat())
r.Brain.Set("foo", "bar")

var out string
r.Brain.Get("foo", &out)
fmt.Println(out)
// Output: "bar"
```

### HTTP

`botopolis` also has an HTTP server built in. One great usecase for this is
webhooks. It exposes [`gorilla/mux`](https://github.com/gorilla/mux) for routing
requests.

```go
r := bot.New(mock.NewChat())
r.Router.HandleFunc("/webhook/name", func(w http.ResponseWriter, r *http.Request) {
		r.Send(bot.Message{Room: "general", Text: "Webhook incoming!"})
		w.Write([]byte("OK"))
	},
).Methods("POST")
```

### Plugins

`botopolis` allows for injection of plugins via `bot.New(chat, ...bot.Plugin)`.
It implements a [`Plugin` interface](https://godoc.org/github.com/botopolis/bot#Plugin)
that allows plugin writers to make use of the full `bot.Robot` runtime on load.

```go
type ExamplePlugin struct { PathName string }

// Load conforms to Plugin interface
func (p ExamplePlugin) Load(r *bot.Robot) {
	r.Router.HandleFunc(e.PathName, func(w http.ResponseWriter, r *http.Request) {
		// special sauce
	}).Methods("GET")
}

// Unload can optionally be implemented for graceful shutdown
func (p ExamplePlugin) Unload(r *bot.Robot) {
	// shutting down
}
```

Here are a couple examples of working plugins:

- [https://github.com/botopolis/bot/help](https://github.com/botopolis/bot/tree/master/help)
- https://github.com/botopolis/oauth2
- https://github.com/botopolis/redis (implements `Unload`)


## Configuration

There are very few configuration options in `botopolis` itself, as it relies on
the introduction of plugins.

### Server port

You can set the server port with an environment variable `PORT`:

```sh
PORT=4567 ./botopolis
```

### Custom logging

You can set up your own logger as long as it satisfies the
[Logger](https://godoc.org/github.com/botopolis/bot#Logger) interface.

```go
robot := bot.New(mychat.Plugin{})
robot.Logger = MyCustomLogger{}
robot.Run()
```
