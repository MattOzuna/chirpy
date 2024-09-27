package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/database"
)

func (cfg apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	body := Chirp{}
	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")

	err := decoder.Decode(&body)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	if !body.isValid() {
		log.Printf("Chirp invalid")
		w.WriteHeader(400)
		message := `{"error": "Chirp is too long"}`
		w.Write([]byte(message))
		return
	}

	dbReq := database.CreateChirpParams{
		Body:   body.Body,
		UserID: body.UserID,
	}

	res, err := cfg.db.CreateChirp(r.Context(), dbReq)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	chirp := Chirp{
		ID:        res.ID,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Body:      res.Body,
		UserID:    res.UserID,
	}
	chirp.profanityFilter()

	dat, err := json.Marshal(&chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(201)
	w.Write(dat)
}
