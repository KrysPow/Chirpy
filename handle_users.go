package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Kryspow/chirpy/internal/auth"
	"github.com/Kryspow/chirpy/internal/database"
	"github.com/google/uuid"
)

type respUser struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (apiC *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type email struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(req.Body)
	posted_email := email{}
	err := decoder.Decode(&posted_email)
	if err != nil {
		fmt.Println("Decoding went wrong", err)
		return
	}

	hashed_pwd, err := auth.HashPassword(posted_email.Password)
	if err != nil {
		fmt.Println("Hashing password went wrong: ", err)
	}

	user, err := apiC.dbQueries.CreateUser(context.Background(), database.CreateUserParams{Email: posted_email.Email,
		HashedPassword: hashed_pwd})
	if err != nil {
		fmt.Println("User creation went wrong: ", err)
		return
	}

	respondWithJson(w, 201, respUser{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (apiC *apiConfig) handlerUpdateUsers(w http.ResponseWriter, req *http.Request) {
	type newPwdEmail struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	new_pwd_email := newPwdEmail{}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "token is missing")
		return
	}

	userId, err := auth.ValidateJWT(token, apiC.secret)
	if err != nil {
		log.Println(err)
		respondWithError(w, 401, "token is not valid")
		return
	}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&new_pwd_email)
	if err != nil {
		log.Println(err)
		return
	}

	hashed_pwd, err := auth.HashPassword(new_pwd_email.Password)
	if err != nil {
		log.Println(err)
		return
	}

	err = apiC.dbQueries.UpdateEmailAndPassword(context.Background(), database.UpdateEmailAndPasswordParams{
		ID:             userId,
		Email:          new_pwd_email.Email,
		HashedPassword: hashed_pwd,
	})
	if err != nil {
		log.Println(err)
		return
	}

	user, err := apiC.dbQueries.GetUserById(context.Background(), userId)
	if err != nil {
		log.Println(err)
		return
	}

	respondWithJson(w, 200, respUser{ID: user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		Token:       token,
		IsChirpyRed: user.IsChirpyRed})

}
