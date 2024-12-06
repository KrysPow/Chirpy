package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Kryspow/chirpy/internal/auth"
)

func (apiC *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Println("Problems getting token from header: ", err)
		return
	}

	refreshToken, err := apiC.dbQueries.GetRefreshToken(context.Background(), token)
	if err != nil {
		respondWithJson(w, 401, "Token does not exist")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithJson(w, 401, "Token expired")
	}

	if refreshToken.RevokedAt.Valid {
		respondWithJson(w, 401, "Token revoked")
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, apiC.secret, time.Hour)
	if err != nil {
		fmt.Println("Error generating JWT: ", err)
		return
	}

	respondWithJson(w, 200, struct {
		Token string `json:"token"`
	}{Token: accessToken})
}

func (apiC *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Println("Problems getting token from header: ", err)
		return
	}

	err = apiC.dbQueries.RevokeToken(context.Background(), token)
	if err != nil {
		fmt.Println("Updating DB went wrong: ", err)
	}
	respondWithJson(w, 204, "")
}
