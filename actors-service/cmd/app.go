package main

import (
	"actors-service/internal/handler"
	"actors-service/internal/postgres"
	"actors-service/internal/repository"
	"actors-service/internal/service"
	"context"
	"log"
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

	actorRepo := repository.NewActorRepository(db)
	actorService := service.NewActorService(actorRepo)

	handler.NewActorHandler(router, actorService)

	log.Println("Actor repository initialized:", actorRepo)
	log.Println("Actor service initialized:", actorService)

	server := http.Server{
		Addr:    ":8003",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Actors microservice start on port 8003")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down actors microservice...")

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

	log.Println("Actor microservice stopped gracefully")

	return nil

}

func main() {
	err := Run()
	if err != nil {
		log.Fatalf("Actor microservice failed: %v", err)
	}
}
