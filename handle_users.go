package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (apiC *apiConfig) handlerUsers(w http.ResponseWriter, req *http.Request) {
	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	posted_email := email{}
	err := decoder.Decode(&posted_email)
	if err != nil {
		fmt.Println("Decoding went wrong", err)
		return
	}

	user, err := apiC.dbQueries.CreateUser(context.Background(), posted_email.Email)
	if err != nil {
		fmt.Println("User creation went wrong: ", err)
		return
	}

	type respUser struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respondWithJson(w, 201, respUser{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
