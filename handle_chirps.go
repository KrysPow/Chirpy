package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Kryspow/chirpy/internal/auth"
	"github.com/Kryspow/chirpy/internal/database"
	"github.com/google/uuid"
)

type respChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (apiC *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, req *http.Request) {
	chirpID := req.PathValue("chirpID")
	uuidChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("%s is not a valid chirpID", chirpID))
		return
	}

	chirp, err := apiC.dbQueries.GetChirp(context.Background(), uuidChirpID)
	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("%s is not a valid chirpID", chirpID))
		return
	}

	respondWithJson(w, 200, respChirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (apiC *apiConfig) handlerGetAllChirps(w http.ResponseWriter, req *http.Request) {
	responseChirps := []respChirp{}

	chirps, err := apiC.dbQueries.GetAllChirps(context.Background())
	if err != nil {
		fmt.Println("Problems getting chirps from dB: ", err)
	}

	for i := range chirps {
		responseChirps = append(responseChirps, respChirp{
			ID:        chirps[i].ID,
			CreatedAt: chirps[i].CreatedAt,
			UpdatedAt: chirps[i].UpdatedAt,
			Body:      chirps[i].Body,
			UserID:    chirps[i].UserID,
		})
	}

	respondWithJson(w, 200, responseChirps)
}

func (apiC *apiConfig) handlerPostChirps(w http.ResponseWriter, req *http.Request) {
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

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Println("Not autorized, ", err)
	}

	user_id, err := auth.ValidateJWT(token, apiC.secret)
	if err != nil {
		respondWithJson(w, 401, "Unauthorized")
		return
	}

	posted_chirp.UserID = user_id

	clean_chirp, err := validateChirp(posted_chirp.Body)
	if err != nil {
		respondWithError(w, 500, "Chirp is too long")
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
