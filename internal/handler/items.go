package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	itemID := r.PathValue("itemID")
	if itemID == "" {
		http.Error(w, "missing 'id' parameter", http.StatusBadRequest)
		return
	}
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		http.Error(w, "cannot parse 'id' parameter", http.StatusInternalServerError)
		return
	}
	data, err := h.Repo.GetItemByID(r.Context(), itemUUID)
	if err != nil {
		msg := fmt.Sprintf("failed fetching item: %v", err)
		http.Error(w, msg, http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
