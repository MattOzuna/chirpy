package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MattOzuna/chirpy/internal/auth"
)

func (cfg apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	body := UserReq{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	user, err := cfg.db.GetUser(r.Context(), body.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	if err := auth.CheckPassword(body.Password, user.HashedPassword); err != nil {
		log.Printf("Password did not match: %s", err)
		w.WriteHeader(401)
		message := `{"error": "Incorrect password"}`
		w.Write([]byte(message))
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, (2 * time.Hour))
	if err != nil {
		log.Printf("Error making JWT: %v", err)
		w.WriteHeader(401)
		message := `{"error":"error making JWT"}`
		w.Write([]byte(message))
		return
	}

	data := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
		Token:     token,
	}

	res, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(res)
}
