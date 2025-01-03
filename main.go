package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MattOzuna/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error connecting to DB: %s", err)
	}

	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"

	cfg := &apiConfig{
		fileserverHits: 0,
		db:             dbQueries,
		platform:       platform,
		secret:         tokenSecret,
	}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot)))

	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.showMetrics)
	mux.HandleFunc("GET /api/chirps", cfg.getAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	mux.HandleFunc("POST /admin/reset", cfg.reset)
	mux.HandleFunc("POST /api/chirps", cfg.createChirp)
	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("POST /api/login", cfg.Login)
	mux.HandleFunc("POST /api/refresh", cfg.Refresh)
	mux.HandleFunc("POST /api/revoke", cfg.Revoke)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.UpgradeUserToRed)
	mux.HandleFunc("PUT /api/users", cfg.editUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	serv := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(serv.ListenAndServe())
}
