package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/get", h.Get)
	mux.HandleFunc("/put", h.Put)
	mux.HandleFunc("/delete", h.Delete)
}
