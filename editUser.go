package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/auth"
	"github.com/MattOzuna/chirpy/internal/database"
)

func (cfg apiConfig) editUser(w http.ResponseWriter, r *http.Request) {

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
	// Get data from body

	body := UserReq{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&body); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	user := database.User{}

	if body.Password != "" && body.Email != "" {
		hashedPassword, err := auth.HashPassword(body.Password)
		if err != nil {
			log.Printf("Error hashing password: %s", err)
			w.WriteHeader(500)
			message := `{"error": "Something went wrong"}`
			w.Write([]byte(message))
			return
		}

		dbUser := database.EditUserParams{
			Email:          body.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		}

		user, err = cfg.db.EditUser(r.Context(), dbUser)
		if err != nil {
			log.Printf("Error editing user: %s", err)
			w.WriteHeader(500)
			message := `{"error": "Something went wrong"}`
			w.Write([]byte(message))
			return
		}
	}

	// if only password was submited
	if body.Email == "" {
		hashedPassword, err := auth.HashPassword(body.Password)
		if err != nil {
			log.Printf("Error hashing password: %s", err)
			w.WriteHeader(500)
			message := `{"error": "Something went wrong"}`
			w.Write([]byte(message))
			return
		}

		dbUser := database.EditUserPasswordParams{
			HashedPassword: hashedPassword,
		}

		user, err = cfg.db.EditUserPassword(r.Context(), dbUser)
		if err != nil {
			log.Printf("Error editing user: %s", err)
			w.WriteHeader(500)
			message := `{"error": "Something went wrong"}`
			w.Write([]byte(message))
			return
		}
	}

	// if only email was submited
	if body.Password == "" {
		dbUser := database.EditUserEmailParams{
			Email: body.Email,
		}

		user, err = cfg.db.EditUserEmail(r.Context(), dbUser)
		if err != nil {
			log.Printf("Error editing user: %s", err)
			w.WriteHeader(500)
			message := `{"error": "Something went wrong"}`
			w.Write([]byte(message))
			return
		}
	}

	//================================================================================//
	// marshall the edited user into JSON and return with 200 code

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
	w.WriteHeader(200)
	w.Write(res)

}
