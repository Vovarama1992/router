package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"router/internal/domain"
)

type Handlers struct {
	svc *domain.Service
}

func NewHandlers(svc *domain.Service) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) ListPeers(w http.ResponseWriter, r *http.Request) {
	peers, err := h.svc.ListPeers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(peers)
}

func (h *Handlers) DisablePeer(w http.ResponseWriter, r *http.Request) {
	tgIDStr := r.URL.Query().Get("telegram_id")
	if tgIDStr == "" {
		http.Error(w, "telegram_id required", http.StatusBadRequest)
		return
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid telegram_id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DisableByTelegramID(r.Context(), tgID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
