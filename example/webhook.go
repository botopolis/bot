package main

import (
	"net/http"

	"github.com/botopolis/bot"
)

type webhook struct{}

func (w webhook) Load(r *bot.Robot) {
	r.Router.HandleFunc(
		"/webhook/name",
		func(w http.ResponseWriter, r *http.Request) {
			// handle webhook
			w.Write([]byte("OK"))
		},
	).Methods("POST")
}
