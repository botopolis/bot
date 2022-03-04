package bot

import (
	"context"
	"net/http"
	"time"
)

type server struct{ http.Server }

func newServer(address string) *server {
	if address == "" || address == ":" {
		address = ":9090"
	}
	return &server{Server: http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}}
}

func (h *server) Load(r *Robot) {
	h.Handler = r.Router
	go func() {
		if err := h.ListenAndServe(); err != nil {
			r.Logger.Errorf("Server error: %s", err.Error())
		}
	}()
}

func (h *server) Unload(r *Robot) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	r.Logger.Info("Shutting down web server")
	if err := h.Shutdown(ctx); err != nil {
		r.Logger.Errorf("Error shutting down: %s\n", err.Error())
	}
}
