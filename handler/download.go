package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	data, err := h.backup.DownloadBackup(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to download backup"})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, data)
}
