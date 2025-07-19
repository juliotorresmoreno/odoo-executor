package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	backups, err := h.backup.ListBackups()
	if err != nil {
		http.Error(w, "Failed to list backups", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(backups)
}
