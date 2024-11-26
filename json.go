package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, msg string) {
	log.Printf("Respoding with error %d: %s", statusCode, msg)

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJson(w, statusCode, errResponse{Error: msg})
}

func respondWithJson(w http.ResponseWriter, statusCode int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error occured during encoding json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Context-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
