package delivery

import (
	"net/http"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(h *Handlers) *Router {
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/peers", h.ListPeers)
	mux.HandleFunc("/api/disable", h.DisablePeer)

	// фронт
	mux.Handle("/", http.FileServer(http.Dir("./front-dist")))

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
