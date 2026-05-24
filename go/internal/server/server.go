package server

import (
	"log"
	"net/http"

	"github.com/bukunya/intero-go/internal/api"
	"github.com/bukunya/intero-go/internal/config"
	"github.com/bukunya/intero-go/internal/satusehat"
)

func NewServer(cfg *config.Config) *http.Server {
	ssClient := satusehat.NewClient(cfg)
	handlers := api.NewHandlers(ssClient)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/patients", handlers.GetPatient)
	mux.HandleFunc("GET /api/practitioners", handlers.GetPractitioner)
	mux.HandleFunc("POST /api/patients", handlers.CreatePatient)
	mux.HandleFunc("POST /api/token", handlers.GetToken)

	handler := loggingMiddleware(mux)

	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
