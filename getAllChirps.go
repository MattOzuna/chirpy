package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/database"
)

func (cfg apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	reMappedChirps := mapChirps(chirps)

	res, err := json.Marshal(&reMappedChirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(res)
}

func mapChirps(c []database.Chirp) []Chirp {
	reMappedChirps := []Chirp{}
	for _, chirp := range c {
		c := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		reMappedChirps = append(reMappedChirps, c)
	}

	return reMappedChirps
}
