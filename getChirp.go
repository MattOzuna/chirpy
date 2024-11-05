package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	reqID := r.PathValue("chirpID")
	i, err := uuid.Parse(reqID)

	if err != nil {
		log.Printf("Error converting id into uuid: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	dbChirp, err := cfg.db.GetChirp(r.Context(), i)
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		w.WriteHeader(404)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	res, err := json.Marshal(&chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
