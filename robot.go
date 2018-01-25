package gobot

import (
	"os"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

// Robot is the central structure for gobot
type Robot struct {
	// Chat adapter
	Chat Chat
	// Data store
	Brain *Brain
	// HTTP Router
	Router *mux.Router
	// Logger with levels
	Logger *logging.Logger

	plugins   *pluginRegistry
	internals *pluginRegistry
	queue     *responderQueue
}

// New creates an instance of a gobot.Robot and loads in the chat adapter
// Typically you would install plugins before running the robot.
func New(c Chat) *Robot {
	r := &Robot{
		Chat:   c,
		Brain:  newBrain(),
		Router: mux.NewRouter(),
		Logger: newLogger(),

		plugins:   newPluginRegistry(),
		internals: newPluginRegistry(),
		queue:     newResponderQueue(10),
	}
	r.internals.Add(
		c,
		newServer(":"+os.Getenv("GOBOT_PORT")),
	)
	return r
}

// Install runs plugins on Robot
func (r *Robot) Install(plugins ...Plugin) {
	r.plugins.Add(plugins...)
}

// Plugin allows you to fetch a loaded plugin for direct use
// See ExamplePluginUsage
func (r *Robot) Plugin(p Plugin) bool {
	return r.plugins.Get(p)
}

// Run is responsible for:
// 1. Loading internals (chat, http)
// 2. Loading all plugins
// 3. Running RTM Slack until termination
// 4. Unloading all plugins in reverse order
// 5. Unloading internals in reverse order
func (r *Robot) Run() {
	r.internals.Load(r)
	r.plugins.Load(r)
	r.queue.Forward(r, r.Chat.Messages())
	r.plugins.Unload(r)
	r.internals.Unload(r)
	os.Exit(0)
}

func (r *Robot) onMessage(t messageType, m Matcher, h hook) {
	if m == nil || h == nil {
		return
	}
	for rs := range r.queue.On(t) {
		if m(&rs) {
			if err := h(rs); err != nil {
				r.Logger.Errorf("Hook error: %s", err.Error())
			}
		}
	}
}

// Hear is triggered on any message event.
func (r *Robot) Hear(m Matcher, h hook) { go r.onMessage(DefaultMessage, m, h) }

// Respond is triggered on messages to the bot
func (r *Robot) Respond(m Matcher, h hook) { go r.onMessage(Response, m, h) }

// Enter is triggered when someone enters a room
func (r *Robot) Enter(h hook) { go r.onMessage(Enter, nil, h) }

// Leave is triggered when someone leaves a room
func (r *Robot) Leave(h hook) { go r.onMessage(Leave, nil, h) }

// Topic is triggered when someone changes the topic
func (r *Robot) Topic(h hook) { go r.onMessage(Topic, nil, h) }

// Username provides the robot's username
func (r *Robot) Username() string { return r.Chat.Username() }

// Debug sets the log-level to debug
func (r *Robot) Debug(debug bool) {
	if debug {
		stdout.SetLevel(logging.DEBUG, "")
	} else {
		stdout.SetLevel(logging.INFO, "")
	}
}
