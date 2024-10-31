package main

import (
	"log"
	"net/http"

	"github.com/MattOzuna/chirpy/internal/auth"
)

func (cfg apiConfig) Revoke(w http.ResponseWriter, r *http.Request) {
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
	// attempt to update revoked_at col in db

	if err := cfg.db.RevokeRefreshToken(r.Context(), authHeader); err != nil {
		log.Printf("Error revoking token: %s", err)
		w.WriteHeader(401)
		message := `{"error": "Refresh Token no found"}`
		w.Write([]byte(message))
		return
	}

	//======================================================================//

	w.WriteHeader(204)
}
