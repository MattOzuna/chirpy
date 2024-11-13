package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if sort == "" || sort == "asc" {
		if s == "" {
			chirps, err = cfg.db.GetAllChirpsAsc(r.Context())
			if err != nil {
				log.Printf("Error getting chirps: %s", err)
				w.WriteHeader(401)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}
		} else {
			user_id, err := uuid.Parse(s)
			if err != nil {
				log.Printf("Error converting id into uuid: %s", err)
				w.WriteHeader(500)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}

			chirps, err = cfg.db.GetAllChirpsByUserIDAsc(r.Context(), user_id)
			if err != nil {
				log.Printf("Error getting chirps: %s", err)
				w.WriteHeader(401)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}
		}
	} else {
		if s == "" {
			chirps, err = cfg.db.GetAllChirpsDesc(r.Context())
			if err != nil {
				log.Printf("Error getting chirps: %s", err)
				w.WriteHeader(401)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}
		} else {
			user_id, err := uuid.Parse(s)
			if err != nil {
				log.Printf("Error converting id into uuid: %s", err)
				w.WriteHeader(500)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}

			chirps, err = cfg.db.GetAllChirpsByUserIDDesc(r.Context(), user_id)
			if err != nil {
				log.Printf("Error getting chirps: %s", err)
				w.WriteHeader(401)
				message := `{"error": "Something went wrong"}`
				w.Write([]byte(message))
				return
			}
		}
	}

	reMappedChirps := mapChirps(chirps)

	res, err := json.Marshal(&reMappedChirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
