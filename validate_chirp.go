package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handlerValidateChip(w http.ResponseWriter, req *http.Request) {
	type Params struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		Valid bool   `json:"valid"`
		Error string `json:"error"`
	}

	decoder := json.NewDecoder(req.Body)
	params := Params{}

	err := decoder.Decode(&params)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{error: %s}", "Json decoding went wrong")))
		w.WriteHeader(500)
		return
	}

	respBody := returnVal{}
	if len(params.Body) > 140 {
		respBody.Error = "Chirp is too long"
		w.WriteHeader(400)
	} else {
		respBody.Valid = true
		w.WriteHeader(200)
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("{error: %s}", "Json encoding went wrong")))
		w.WriteHeader(500)
		return
	}

	w.Write(dat)

}
