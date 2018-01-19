package gobot

import (
	"net/http"
	"os"

	"github.com/berfarah/gobot/brain"
	"github.com/gorilla/mux"
	"github.com/op/go-logging"
)

type Plugin func(*Robot)

type Robot struct {
	server *http.Server

	Adapter Adapter
	Brain   *brain.Brain
	Router  *mux.Router
	Logger  *logging.Logger

	// Bot info
	Name string

	// Event management
	triggers map[string][]*Hook
	events   *EventDispatcher
}

func New(secret string, a Adapter, store brain.Store) *Robot {
	return &Robot{
		Adapter: a,
		Brain:   brain.New(store),
		Router:  mux.NewRouter(),
		Logger:  newLogger(),

		triggers: newTriggers(),
		events:   NewEventDispatcher(),
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
	// r.listenEvents()
	r.stopHTTP()
	os.Exit(0)
}
