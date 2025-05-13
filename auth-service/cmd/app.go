package main

import (
	"auth-service/internal/handlers"
	"auth-service/internal/postgres"
	"auth-service/internal/repository"
	"auth-service/internal/service"
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

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository)

	handlers.NewAuthHandler(router, authService)

	log.Println("Auth repository initialized", authRepository)
	log.Println("Auth service initialized", authService)

	server := http.Server{
		Addr:    ":8001",
		Handler: router,
	}

	log.Println("Auth service start on 8001 port")

	server.ListenAndServe()

	return nil
}

func main() {

	err := Run()

	if err != nil {
		log.Fatalf("Auth microservice failed: %v", err)
	}
}
