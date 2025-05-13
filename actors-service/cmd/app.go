package main

import (
	"actors-service/internal/postgres"
	"actors-service/internal/repository"
	"log"
	"net/http"
)

func Run() error {
	router := http.NewServeMux()
	db, err := postgres.NewConnectDb()
	if err != nil {
		log.Panicln(err)
	}

	log.Println("Postgres DB connected: ", db)

	defer db.Close()

	actorRepo := repository.NewActorRepository(db)

	log.Println("Actor repository initialized:", actorRepo)

	server := http.Server{
		Addr:    ":8003",
		Handler: router,
	}
	log.Println("Actors microservice start on port 8003")

	server.ListenAndServe()

	return nil

}

func main() {
	err := Run()
	if err != nil {
		log.Fatalf("Actors microservice failed: %v", err)
	}
}
