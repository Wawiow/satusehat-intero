package main

import (
	"log"

	"github.com/bukunya/intero-go/internal/config"
	"github.com/bukunya/intero-go/internal/server"
)

func main() {
	log.Println("Starting server...")

	cfg := config.LoadConfig()

	srv := server.NewServer(cfg)

	log.Printf("Server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
