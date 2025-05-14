package main

import (
	"context"
	"log"
	"movies-service/internal/handler"
	"movies-service/internal/postgres"
	"movies-service/internal/repository"
	"movies-service/internal/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	movieService := service.NewMovieService(movieRepository)

	handler.NewMovieHandler(router, movieService)

	log.Println("Movie repository initialized:", movieRepository)
	log.Println("Movie service initialized:", movieService)

	server := http.Server{
		Addr:    ":8002",
		Handler: router,
	}
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Movies microservice start on port 8002")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit

	log.Println("Shutting down movie microservice...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	err = db.Close()

	if err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Movie microservice stopped gracefully")

	return nil

}

func main() {
	err := Run()
	if err != nil {
		log.Fatalf("Movie microservice failed: %v", err)
	}
}
