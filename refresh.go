package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/MattOzuna/chirpy/internal/auth"
)

type RefreshRes struct {
	Token string `json:"token"`
}

func (cfg apiConfig) Refresh(w http.ResponseWriter, r *http.Request) {

	//======================================================================//
	//Getting Refresh token from auth header
	//expect header to look like Authorization: Bearer <token>

	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting auth header: %s", err)
		w.WriteHeader(500)
		message := `{"error": "Something went wrong"}`
		w.Write([]byte(message))
		return
	}

	//======================================================================//
	//Check to see if the refresh token in in the DB

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), authHeader)
	if err != nil || refreshToken.RevokedAt.Valid {
		log.Printf("Error getting token: %s", err)
		w.WriteHeader(401)
		message := `{"error": "Refresh Token not found or has been revoked"}`
		w.Write([]byte(message))
		return
	}

	//======================================================================//
	// if refreshToken exists in DB, make a JWT to retrun

	token, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, (1 * time.Hour))
	if err != nil {
		log.Printf("Error making JWT: %v", err)
		w.WriteHeader(401)
		message := `{"error":"error making JWT"}`
		w.Write([]byte(message))
		return
	}

	//======================================================================//
	//marshal the response into a RefreshRes struct

	data := RefreshRes{
		Token: token,
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
