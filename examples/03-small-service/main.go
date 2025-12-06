// Package main
package main

import (
	"log"
	"net/http"

	"small-service/internal/config"
	"small-service/internal/user"
)

func main() {

	cfg := config.Load()
	log.Printf("Starting server with config: %+v", cfg)

	storage := user.NewStorage()
	storage.Create("Alice", "alice@example.com")
	storage.Create("Bob", "bob@example.com")

	handler := user.NewHandler(storage)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	addr := cfg.ServerAddr
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
