package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

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
	_ = json.NewEncoder(w).Encode(peers)
}

func (h *Handlers) SetUserUntil(w http.ResponseWriter, r *http.Request) {
	tgIDStr := r.URL.Query().Get("telegram_id")
	untilStr := r.URL.Query().Get("until")

	if tgIDStr == "" || untilStr == "" {
		http.Error(w, "telegram_id and until required", http.StatusBadRequest)
		return
	}

	tgID, err := strconv.ParseInt(tgIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid telegram_id", http.StatusBadRequest)
		return
	}

	until, err := time.Parse(time.RFC3339, untilStr)
	if err != nil {
		http.Error(w, "invalid until format (use RFC3339)", http.StatusBadRequest)
		return
	}

	if err := h.svc.SetUserUntil(r.Context(), tgID, until); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
}
