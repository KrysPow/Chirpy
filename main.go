package main

import (
	"fmt"
	"net/http"
)

func main() {
	servMux := http.NewServeMux()
	server := http.Server{
		Handler: servMux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	fmt.Println(err)
}
