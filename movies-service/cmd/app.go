package main

import (
	"log"
	"movies-service/internal/postgres"
	"movies-service/internal/repository"
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

	movieRepository := repository.NewMovieRepository(db)

	log.Println("Movie repository initialized:", movieRepository)

	server := http.Server{
		Addr:    ":8002",
		Handler: router,
	}
	log.Println("Movies microservice start on port 8002")

	server.ListenAndServe()

	return nil

}

func main() {
	err := Run()
	if err != nil {
		log.Fatalf("Movies microservice failed: %v", err)
	}
}
