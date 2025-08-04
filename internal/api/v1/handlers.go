package v1

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Handler struct{}

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	payload := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		slog.Default().ErrorContext(r.Context(), "encoding payload", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
