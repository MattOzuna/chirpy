package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MattOzuna/chirpy/internal/auth"
	"github.com/MattOzuna/chirpy/internal/database"
	"github.com/google/uuid"
)

type UserReq struct {
	Password  string        `json:"password"`
	Email     string        `json:"email"`
	ExpiresIn time.Duration `json:"expires_in_seconds"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	body := UserReq{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	if body.Password == "" {
		log.Print("Password not sent with request")
		w.WriteHeader(500)
		message := `{"error": "Must include pasword"}`
		w.Write([]byte(message))
		return
	}

	hashedPassword, err := auth.HashPassword(body.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	dbUser := database.CreateUserParams{
		Email:          body.Email,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), dbUser)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	data := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Email:     user.Email,
	}

	res, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(res)
}
