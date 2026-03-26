package api

import (
	"encoding/json"
	"net/http"

	"github.com/OLABALADE/FluxKV/internal/store"
)

type Handler struct {
	store store.Store
}

func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

type Request struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	req := &Request{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if err := h.store.Put(req.Key, req.Value); err != nil {
		http.Error(w, "Failed to add Key to store", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val, err := h.store.Get(key)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"value": val,
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if err := h.store.Delete(key); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
