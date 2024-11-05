package main

import (
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	//================================================================================//
	//get auth header and validate JWT

	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting auth header: %s", err)
		w.WriteHeader(401)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	userID, err := auth.ValidateJWT(authHeader, cfg.secret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		w.WriteHeader(401)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}
	//================================================================================//
	// Get chirp ID from params and check userID from chirp matches auth header

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

	if userID != dbChirp.UserID {
		log.Println("UserID header and chirp UserID do not match")
		w.WriteHeader(403)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	//================================================================================//
	// Delete chirp from DB
	err = cfg.db.DeleteChirp(r.Context(), i)
	if err != nil {
		log.Printf("Error deleting chirps: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	w.WriteHeader(204)
}
