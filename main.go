package main

import (
	"fmt"
	"net/http"
)

func readiness(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{Content-Type: text/plain; charset=utf-8, Status: OK}"))
}

func main() {
	servMux := http.NewServeMux()
	server := http.Server{
		Handler: servMux,
		Addr:    ":8080",
	}

	servMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	servMux.HandleFunc("/healthz", readiness)

	err := server.ListenAndServe()
	fmt.Println(err)
}
