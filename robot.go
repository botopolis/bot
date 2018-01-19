package gobot

import (
	"net/http"
	"os"

	"github.com/berfarah/gobot/brain"
	mux "github.com/gorilla/mux"
	"github.com/nlopes/slack"
	"github.com/op/go-logging"
)

type Plugin func(*Robot)

type Robot struct {
	rtm    *slack.RTM
	api    *slack.Client
	server *http.Server

	Brain  *brain.Brain
	Router *mux.Router
	Logger *logging.Logger

	// Bot info
	ID   string
	Name string

	// Event management
	triggers map[string]trigger
	events   *EventDispatcher
}

func New(secret string, store brain.Store) *Robot {
	return &Robot{
		api:      slack.New(secret),
		triggers: newTriggers(),
		events:   NewEventDispatcher(),
		Brain:    brain.New(store),
		Router:   mux.NewRouter(),
		Logger:   newLogger(),
	}
}

// Install runs plugins on Robot
func (r *Robot) Install(plugins ...Plugin) {
	for _, p := range plugins {
		if p != nil {
			p(r)
		}
	}
	r.Logger.Debugf("%d plugins loaded", len(plugins))
}

func (r *Robot) Start(address string) {
	r.listenHTTP(address)
	r.listenEvents()
	r.stopHTTP()
	os.Exit(0)
}
