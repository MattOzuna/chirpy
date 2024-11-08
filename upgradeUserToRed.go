package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type UserUpgrade struct {
	Event string `json:"event"`
	Data  UserID `json:"data"`
}

type UserID struct {
	UserID uuid.UUID `json:"user_id"`
}

func (cfg apiConfig) UpgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	body := UserUpgrade{}

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	if body.Event != "user.upgraded" {
		log.Println("Event not 'user.upgraded'")
		w.WriteHeader(204)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	err := cfg.db.UpgradeUserToRed(r.Context(), body.Data.UserID)
	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		w.WriteHeader(404)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	w.WriteHeader(204)

}
