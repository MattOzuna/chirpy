package main

import (
	"strings"
)

type Chirp struct {
	Body        string `json:"body"`
	CleanedBody string `json:"cleaned_body"`
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
