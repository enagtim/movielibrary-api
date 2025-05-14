package main

import (
	"api-gateway/middleware"
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func proxyToService(target string, prefix string) http.Handler {
	return http.StripPrefix(prefix, httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	}))
}

func main() {
	os.Setenv("JWT_SECRET", os.Getenv("JWT_SECRET"))

	// Регистрация и авторизация

	http.Handle("/api/auth/", proxyToService("auth:8001", "/api/auth"))

	// actors service для пользователя

	http.Handle("/api/actors", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("actors:8003", "/api"),
	))
	http.Handle("/api/actors/", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("actors:8003", "/api"),
	))

	// actors service для админа

	http.Handle("/api/admin/actors", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("actors:8003", "/api/admin"),
	))

	http.Handle("/api/admin/actors/", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("actors:8003", "/api/admin"),
	))

	// movies service для пользователя

	http.Handle("/api/movies", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("movies:8002", "/api"),
	))

	http.Handle("/api/movies/", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("movies:8002", "/api"),
	))

	// movies service для админа

	http.Handle("/api/admin/movies", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("movies:8002", "/api/admin"),
	))

	http.Handle("/api/admin/movies/", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("movies:8002", "/api/admin"),
	))

	server := http.Server{
		Addr: ":8080",
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("API Gateway started on :8080")
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := server.Shutdown(ctx)

	if err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("API-gateway stopped gracefully")

}
