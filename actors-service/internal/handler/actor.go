package handler

import (
	"actors-service/internal/payload"
	"actors-service/internal/service"
	"actors-service/pkg/req"
	"actors-service/pkg/res"
	"context"
	"net/http"
	"strconv"
	"time"
)

type ActorHandler struct {
	ActorService *service.ActorService
}

func NewActorHandler(router *http.ServeMux, actorService *service.ActorService) {
	handler := &ActorHandler{
		ActorService: actorService,
	}

	router.HandleFunc("POST /actors", handler.CreateActor())
	router.HandleFunc("GET /actors", handler.GetActorsWithMovies())
	router.HandleFunc("GET /actors/{id}", handler.GetActorByID())
	router.HandleFunc("PUT /actors/{id}", handler.FullUpdateActorByID())
	router.HandleFunc("PATCH /actors/{id}", handler.PartialUpdateActorByID())
	router.HandleFunc("DELETE /actors/{id}", handler.DeleteActor())
}

func (h *ActorHandler) CreateActor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		body, err := req.DecodedAndValidatedBody[payload.ActorPayload](r.Body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusBadRequest)
			return
		}

		actorID, err := h.ActorService.Create(ctx, &body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &payload.CreatedActorResponse{
			ActorID: actorID,
			Message: "Actor was created",
		}

		res.ResJson(w, data, http.StatusCreated)
	}
}
func (h *ActorHandler) GetActorsWithMovies() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		actors, err := h.ActorService.GetActorsWithMovies(ctx)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := &payload.GetActorResponse{
			Data: actors,
		}

		res.ResJson(w, data, http.StatusCreated)

	}
}

func (h *ActorHandler) GetActorByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		param := r.PathValue("id")

		id, err := strconv.Atoi(param)

		if err != nil {
			res.ErrResJson(w, "Invalid actor ID", http.StatusBadRequest)
			return
		}

		actor, err := h.ActorService.GetByID(ctx, uint(id))

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.ResJson(w, actor, http.StatusOK)

	}
}

func (h *ActorHandler) FullUpdateActorByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		param := r.PathValue("id")

		id, err := strconv.Atoi(param)

		if err != nil {
			res.ErrResJson(w, "Invalid actor ID", http.StatusBadRequest)
			return
		}
		body, err := req.DecodedAndValidatedBody[payload.ActorPayload](r.Body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.ActorService.FullUpdate(ctx, uint(id), &body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *ActorHandler) PartialUpdateActorByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		param := r.PathValue("id")

		id, err := strconv.Atoi(param)

		if err != nil {
			res.ErrResJson(w, "Invalid actor ID", http.StatusBadRequest)
			return
		}
		body, err := req.DecodedAndValidatedBody[payload.PartialUpdateActorPayload](r.Body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = h.ActorService.PartialUpdate(ctx, uint(id), &body)

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	}
}

func (h *ActorHandler) DeleteActor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)

		defer cancel()

		param := r.PathValue("id")

		id, err := strconv.Atoi(param)

		if err != nil {
			res.ErrResJson(w, "Invalid actor ID", http.StatusBadRequest)
			return
		}

		err = h.ActorService.Delete(ctx, uint(id))

		if err != nil {
			res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
