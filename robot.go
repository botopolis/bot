package bot

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

const timeout = 15

// Robot is the central structure for bot
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

// New creates an instance of a bot.Robot and loads in the chat adapter
// Typically you would install plugins before running the robot.
func New(c Chat, plugins ...Plugin) *Robot {
	r := &Robot{
		Chat:   c,
		Brain:  NewBrain(),
		Router: mux.NewRouter(),
		Logger: newLogger(),

		plugins:   newPluginRegistry(),
		internals: newPluginRegistry(),
		queue:     newResponderQueue(10),
	}
	r.internals.Add(
		c,
		newServer(":"+os.Getenv("PORT")),
	)
	r.plugins.Add(plugins...)
	return r
}

// Plugin allows you to fetch a loaded plugin for direct use
// See ExamplePluginUsage
func (r *Robot) Plugin(p Plugin) bool {
	ok := r.plugins.Get(p)
	if !ok {
		return r.internals.Get(p)
	}
	return ok
}

// Run is responsible for:
// 1. Loading internals (chat, http)
// 2. Loading all plugins
// 3. Running RTM Slack until termination
// 4. Unloading all plugins in reverse order
// 5. Unloading internals in reverse order
func (r *Robot) Run() {
	r.gracefulShutdown()
	r.internals.Load(r)
	r.plugins.Load(r)
	r.queue.Forward(r, r.Chat.Messages())
	r.stop()
}

func (r *Robot) gracefulShutdown() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		done := make(chan bool)
		go func() {
			r.stop()
			done <- true
		}()
		select {
		case <-time.After(timeout * time.Second):
			r.Logger.Info("Force quitting after %ds timeout", timeout)
			os.Exit(1)
		case <-done:
			os.Exit(0)
		}
	}()
}

func (r *Robot) stop() {
	r.Logger.Info("Shutting down...")
	r.plugins.Unload(r)
	r.internals.Unload(r)
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
