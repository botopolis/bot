package gobot

import (
	"context"
	"net/http"
	"time"
)

func (r *Robot) listenHTTP(address string) {
	r.server = &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r.Router,
	}
	go func() {
		if err := r.server.ListenAndServe(); err != nil {
			r.Logger.Errorf("Server error: %s", err.Error())
		}
	}()
}

func (r *Robot) stopHTTP() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	r.Logger.Info("Shutting down web server")
	r.server.Shutdown(ctx)
}
