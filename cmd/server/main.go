package main

import (
	"log"
	"net/http"

	"github.com/OLABALADE/FluxKV/internal/api"
	"github.com/OLABALADE/FluxKV/internal/store"
)

func main() {
	port := "8080"
	mux := http.NewServeMux()
	kvStore := store.NewMemoryStore()
	handler := api.NewHandler(kvStore)
	api.RegisterRoutes(mux, handler)

	log.Printf("FluxKV node starting on port %s...\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
