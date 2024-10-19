package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits int
	db             *database.Queries
	platform       string
	secret         string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) showMetrics(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	content := fmt.Sprintf(`
<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits)
	w.Write([]byte(content))
}

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("reset metrics"))
}
