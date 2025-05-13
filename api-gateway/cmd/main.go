package main

import (
	"api-gateway/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func proxyToService(target string, prefix string) http.Handler {
	return http.StripPrefix(prefix, httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   target,
	}))
}

func main() {
	os.Setenv("JWT_SECRET", os.Getenv("JWT_SECRET"))

	http.Handle("/api/auth/", proxyToService("auth:8001", "/api/auth"))

	http.Handle("/api/movies/", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("movies:8002", "/api/movies"),
	))

	http.Handle("/api/actors/", middleware.CheckRoleAndMethod(
		"user",
		[]string{"GET"},
		proxyToService("actors:8003", "/api/actors"),
	))

	http.Handle("/api/admin/movies/", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("movies:8002", "/api/admin/movies"),
	))

	http.Handle("/api/admin/actors/", middleware.CheckRoleAndMethod(
		"admin",
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		proxyToService("actors:8003", "/api/admin/actors"),
	))

	log.Println("API Gateway started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
