package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {
	chirp := Chirp{}
	decoder := json.NewDecoder(r.Body)
	w.Header().Set("Content-Type", "application/json")

	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	if !chirp.isValid() {
		log.Printf("Chirp invalid")
		w.WriteHeader(400)
		message := `{"error": "Chirp is too long"}`
		w.Write([]byte(message))
		return
	}

	chirp.profanityFilter()
	dat, err := json.Marshal(&chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}
