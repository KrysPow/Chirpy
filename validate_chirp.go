package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func handlerValidateChip(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Body string `json:"body"`
	}

	type cleanBody struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Params{}

	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{error: %s}", "Json decoding went wrong")))
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		respondWithJson(w, 200, cleanBody{CleanBody: censorProfaneWords(params.Body)})

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
