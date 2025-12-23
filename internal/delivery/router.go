package delivery

import (
	"github.com/go-chi/chi/v5"
)

func RegisterVPNRoutes(r chi.Router, h *VPNHandler) {
	r.Post("/vpn/peer", h.CreatePeer)
}
