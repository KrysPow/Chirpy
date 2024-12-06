package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Kryspow/chirpy/internal/auth"
	"github.com/Kryspow/chirpy/internal/database"
)

func (apiC *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req_body := reqBody{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&req_body)
	if err != nil {
		fmt.Println("Decoding json went wrong: ", err)
		return
	}

	user, err := apiC.dbQueries.GetUserByEmail(context.Background(), req_body.Email)
	if err != nil {
		fmt.Println("User does not exist: ", err)
		return
	}

	err = auth.CheckPasswordHash(req_body.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, apiC.secret, time.Hour)

	if err != nil {
		fmt.Println("Problems creating JWT token: ", err)
		return
	}

	fresh_token, err := auth.MakeFreshToken()
	if err != nil {
		fmt.Println("Token couldn't be refreshed, ", err)
		return
	}

	refreshToken, err := apiC.dbQueries.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{Token: fresh_token, UserID: user.ID})
	if err != nil {
		fmt.Println("Writing refresh token to DB failed: ", err)
		return
	}

	respondWithJson(w, 200, respUser{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})

}
