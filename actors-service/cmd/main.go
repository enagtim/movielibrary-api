package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8003",
		Handler: router,
	}
	server.ListenAndServe()
	log.Println("Actors-service start on port 8003")
}
