package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Kryspow/chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiC *apiConfig) handlerChirps(w http.ResponseWriter, req *http.Request) {
	type postChirp struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	posted_chirp := postChirp{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&posted_chirp)
	if err != nil {
		fmt.Println("Decoding went wrong: ", err)
		return
	}

	clean_chirp, err := validateChirp(posted_chirp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	chirp, err := apiC.dbQueries.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   clean_chirp,
		UserID: posted_chirp.UserID,
	})
	if err != nil {
		fmt.Println("Saving chirp to database failed: ", err)
		return
	}

	type respChirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	respondWithJson(w, 201, respChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", fmt.Errorf("chirp is too long, only 140 characters are allowed")
	} else {
		return censorProfaneWords(chirp), nil
	}
}

func censorProfaneWords(s string) string {
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax"}

	var clean_s string
	sSplitted := strings.Split(s, " ")
	for _, profWord := range profaneWords {
		for i := range sSplitted {
			if strings.Contains(strings.ToLower(sSplitted[i]), profWord) {
				sSplitted[i] = strings.Replace(strings.ToLower(sSplitted[i]), profWord, "****", 1)
			}
		}
		clean_s = strings.Join(sSplitted, " ")
	}
	return clean_s
}
