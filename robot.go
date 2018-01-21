package gobot

import (
	"net/http"
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

	plugins *pluginRegistry
	server  *http.Server
	queue   *ResponderQueue
}

// New creates an instance of a gobot.Robot and loads in the chat adapter
// Typically you would install plugins before running the robot.
func New(c Chat) *Robot {
	r := &Robot{
		Chat:   c,
		Brain:  newBrain(),
		Router: mux.NewRouter(),
		Logger: newLogger(),

		plugins: newPluginRegistry(),
		queue:   NewResponderQueue(32),
	}
	c.Load(r)
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
// 1. Launching the HTTP server
// 2. Loading all plugins
// 3. Running RTM Slack until termination
// 4. Unloading all plugins
// 5. Shutting down the HTTP server
func (r *Robot) Run(address string) {
	r.listenHTTP(address)
	r.plugins.Load(r)
	r.queue.Forward(r, r.Chat.Messages())
	r.plugins.Unload(r)
	r.stopHTTP()
	os.Exit(0)
}

func (r *Robot) onMessage(et EventType, h *Hook) {
	for rs := range r.queue.On(et) {
		h.Run(&rs)
	}
}

// Hear is triggered on any message event.
func (r *Robot) Hear(h *Hook) { go r.onMessage(MessageEvent, h) }

// Respond is triggered on messages to the bot
func (r *Robot) Respond(h *Hook) { go r.onMessage(RespondEvent, h) }

// Enter is triggered when someone enters a room
func (r *Robot) Enter(h *Hook) { go r.onMessage(EnterEvent, h) }

// Leave is triggered when someone leaves a room
func (r *Robot) Leave(h *Hook) { go r.onMessage(LeaveEvent, h) }

// Topic is triggered when someone changes the topic
func (r *Robot) Topic(h *Hook) { go r.onMessage(TopicEvent, h) }
