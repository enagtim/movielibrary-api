package handlers

import (
	"auth-service/internal/payload"
	"auth-service/internal/service"
	"auth-service/pkg/req"
	"auth-service/pkg/res"
	"context"
	"net/http"
	"time"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(router *http.ServeMux, authService *service.AuthService) {
	handler := &AuthHandler{
		AuthService: authService,
	}

	router.HandleFunc("POST /auth/register", handler.RegisterUser())

}

func (h *AuthHandler) RegisterUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		body, err := req.DecodedAndValidatedBody[payload.AuthRegisterPayload](r.Body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusBadRequest)
			return
		}

		userID, err := h.AuthService.Register(ctx, &body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &payload.AuthRegisterResponse{
			UserID:  userID,
			Message: "User successfully created",
		}

		res.ResJson(w, data, http.StatusCreated)

	}
}

func (h *AuthHandler) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		body, err := req.DecodedAndValidatedBody[payload.AuthLoginPayload](r.Body)
		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := h.AuthService.Login(ctx, &body)
		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusUnauthorized)
			return
		}

		data := &payload.AuthLoginResponse{
			Token: token,
		}

		res.ResJson(w, data, http.StatusOK)
	}
}
