package delivery

import (
	"net/http"
	"os"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(h *Handlers) *Router {
	mux := http.NewServeMux()

	// API
	mux.HandleFunc("/api/peers", h.ListPeers)
	mux.HandleFunc("/api/user/until", h.SetUserUntil)

	// фронт
	fileServer := http.FileServer(http.Dir("./front-dist"))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := "./front-dist" + r.URL.Path

		if _, err := os.Stat(path); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, "./front-dist/index.html")
	}))

	return &Router{mux: mux}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
