// @title           SatuSehat Intero API
// @version         1.0
// @description     This is a SatuSehat integration and interoperability API client wrapper.
// @host            localhost:8083
// @BasePath        /
package main

import (
	"log"

	"github.com/bukunya/intero-go/internal/config"
	"github.com/bukunya/intero-go/internal/database"
	"github.com/bukunya/intero-go/internal/server"
)

func main() {
	log.Println("Starting server...")

	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	srv := server.NewServer(cfg, db)

	log.Printf("Server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
