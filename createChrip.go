package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/auth"
	"github.com/MattOzuna/chirpy/internal/database"
)

func (cfg apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting auth header: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	UserID, err := auth.ValidateJWT(authHeader, cfg.secret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	body := Chirp{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&body)
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
		UserID: UserID,
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(dat)
}
