package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/postgres"
	"auth-service/internal/repository"
	"auth-service/internal/service"
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

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)

	handlers.NewAuthHandler(router, authService)

	log.Println("Auth repository initialized", authRepository)
	log.Println("Auth service initialized", authService)

	server := http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Auth microservice start on port 8001")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit

	log.Println("Shutting down auth microservice...")

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

	log.Println("Auth microservice stopped gracefully")

	return nil
}

func main() {
	err := Run()
	if err != nil {
		log.Fatalf("Auth microservice failed: %v", err)
	}
}
