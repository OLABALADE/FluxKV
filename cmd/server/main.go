package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/OLABALADE/FluxKV/internal/api"
	"github.com/OLABALADE/FluxKV/internal/cluster"
	"github.com/OLABALADE/FluxKV/internal/store"
)

func main() {
	port := getEnv("PORT", "8080")
	nodeID := getEnv("NODE_ID", port)
	peersEnv := getEnv("PEERS", port)

	var peers []string
	if peersEnv != "" {
		peers = strings.Split(peersEnv, ",")
	}

	mux := http.NewServeMux()
	kvStore := store.NewMemoryStore()

	node := cluster.NewNode(nodeID, "localhost:"+port, peers)
	handler := api.NewHandler(kvStore, node)

	api.RegisterRoutes(mux, handler)

	log.Printf("FluxKV node starting on port %s...\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return env
}
