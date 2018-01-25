package main

import (
	"net/http"

	"github.com/berfarah/gobot"
)

type webhook struct{}

func (w webhook) Load(r *gobot.Robot) {
	r.Router.HandleFunc(
		"/webhook/name",
		func(w http.ResponseWriter, r *http.Request) {
			// handle webhook
			w.Write([]byte("OK"))
		},
	).Methods("POST")
}
