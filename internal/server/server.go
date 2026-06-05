package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bukunya/intero-go/internal/api"
	"github.com/bukunya/intero-go/internal/config"
	"github.com/bukunya/intero-go/internal/satusehat"

	_ "github.com/bukunya/intero-go/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewServer(cfg *config.Config, db *sql.DB) *http.Server {
	ssClient := satusehat.NewClient(cfg)
	handlers := api.NewHandlers(ssClient, db)

	mux := http.NewServeMux()

	// Swagger
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// OK
	mux.HandleFunc("GET /api/patients", handlers.GetPatient)
	mux.HandleFunc("GET /api/local/patients", handlers.GetAllLocalPatients)
	mux.HandleFunc("POST /api/patients", handlers.CreatePatient)
	mux.HandleFunc("POST /api/token", handlers.GetToken)

	// Not OK
	mux.HandleFunc("GET /api/practitioners", handlers.GetPractitioners)
	mux.HandleFunc("GET /api/local/practitioners", handlers.GetAllLocalPractitioners)

	// OK
	mux.HandleFunc("GET /api/locations", handlers.GetLocations)
	mux.HandleFunc("POST /api/locations", handlers.CreateLocation)

	// Not OK

	// OK

	// Not OK
	mux.HandleFunc("GET /api/encounters/{id}", handlers.GetEncounterById)
	mux.HandleFunc("GET /api/local/encounters", handlers.GetAllLocalEncounters)
	mux.HandleFunc("POST /api/encounters", handlers.CreateEncounter)
	mux.HandleFunc("PUT /api/encounters/{id}", handlers.UpdateEncounterStatus)

	handler := loggingMiddleware(mux)

	return &http.Server{
		Addr:    ":8083",
		Handler: handler,
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
