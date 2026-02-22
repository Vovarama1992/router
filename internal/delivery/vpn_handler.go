package delivery

import (
	"net/http"
	"router/internal/domain"
)

type VPNHandler struct {
	svc *domain.Service
}

func NewVPNHandler(svc *domain.Service) *VPNHandler {
	return &VPNHandler{svc: svc}
}

// POST /vpn/peer
func (h *VPNHandler) CreatePeer(w http.ResponseWriter, r *http.Request) {
	peer, err := h.svc.CreatePeer(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(peer.Link))
}
