package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const RoleKey key = "userRole"

func CheckRoleAndMethod(requiredRole string, allowedMethods []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing auth header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || (requiredRole == "admin" && role != "admin") {
			http.Error(w, "access denied", http.StatusForbidden)
			return
		}

		if !isMethodAllowed(r.Method, allowedMethods) {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := context.WithValue(r.Context(), RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func isMethodAllowed(requestMethod string, allowedMethods []string) bool {
	for _, method := range allowedMethods {
		if method == requestMethod {
			return true
		}
	}
	return false
}
