package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Kryspow/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (apiC *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, req *http.Request) {
	type Data struct {
		UserID string `json:"user_id"`
	}

	type Webhook struct {
		Event string `json:"event"`
		Data  Data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		log.Println(err)
	}

	if apiKey != apiC.polkaKey {
		respondWithError(w, 401, "wrong API-key")
		return
	}

	webhook := Webhook{}

	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&webhook)
	if err != nil {
		log.Println(err)
		return
	}

	if webhook.Event != "user.upgraded" {
		respondWithJson(w, 204, "")
		return
	}
	uuid_userId, err := uuid.Parse(webhook.Data.UserID)
	if err != nil {
		log.Println(err)
	}

	err = apiC.dbQueries.UpdateToRedById(context.Background(), uuid_userId)
	if err != nil {
		log.Println(err)
		respondWithError(w, 404, "user can't be found")
		return
	}
	respondWithJson(w, 204, "")

}
