package main

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Body        string    `json:"body"`
	UserID      uuid.UUID `json:"user_id"`
	CleanedBody string    `json:"cleaned_body"`
}

func (c *Chirp) isValid() bool {
	return len(c.Body) <= 140
}

func (c *Chirp) profanityFilter() {
	fields := strings.Fields(c.Body)
	for i, field := range fields {
		field = strings.ToLower(field)
		if field == "kerfuffle" || field == "sharbert" || field == "fornax" {
			fields[i] = "****"
		}
	}
	str := strings.Join(fields, " ")
	c.CleanedBody = str
}
