package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MattOzuna/chirpy/internal/auth"
	"github.com/MattOzuna/chirpy/internal/database"
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, (1 * time.Hour))
	if err != nil {
		log.Printf("Error making JWT: %v", err)
		w.WriteHeader(401)
		message := `{"error":"error making JWT"}`
		w.Write([]byte(message))
		return
	}

	// =============== Make a Refresh Token and store in DB ===================
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error making refresh token: %v", err)
		w.WriteHeader(401)
		message := `{"error":"error making refresh token"}`
		w.Write([]byte(message))
		return
	}

	refreshTokenForDB := database.InsertRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	}

	_, err = cfg.db.InsertRefreshToken(r.Context(), refreshTokenForDB)
	if err != nil {
		log.Printf("Error making refresh token in DB: %v", err)
		w.WriteHeader(401)
		message := `{"error":"error making refresh token in DB"}`
		w.Write([]byte(message))
		return
	}
	//======================================================================//

	data := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt.Time,
		UpdatedAt:    user.UpdatedAt.Time,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed.Bool,
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
