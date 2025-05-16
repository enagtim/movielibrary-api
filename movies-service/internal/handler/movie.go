package handler

import (
	"context"
	"movies-service/internal/payload"
	"movies-service/internal/service"
	"movies-service/pkg/req"
	"movies-service/pkg/res"
	"net/http"
	"strconv"
	"time"
)

type MovieHandler struct {
	MovieService *service.MovieService
}

func NewMovieHandler(router *http.ServeMux, movieService *service.MovieService) {
	handler := &MovieHandler{
		MovieService: movieService,
	}

	router.HandleFunc("POST /movies", handler.CreateMovie)
	router.HandleFunc("GET /movies", handler.GetAllMovies)
	router.HandleFunc("GET /movies/{id}", handler.GetMovieByID)
	router.HandleFunc("PUT /movies/{id}", handler.FullUpdateMovieByID)
	router.HandleFunc("PATCH /movies/{id}", handler.PartialUpdateMovieByID)
	router.HandleFunc("DELETE /movies/{id}", handler.DeleteMovieByID)
	router.HandleFunc("GET /movies/search/title", handler.SearchMovieByTitle)
	router.HandleFunc("GET /movies/search/actorname", handler.SearchMovieByActorName)
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	body, err := req.DecodedAndValidatedBody[payload.MoviePayload](r.Body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	movieID, err := h.MovieService.Create(ctx, &body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.CreateMovieResponse{
		MovieID: movieID,
		Message: "Movie was created",
	}

	res.ResJson(w, data, http.StatusCreated)

}

func (h *MovieHandler) GetAllMovies(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	sortBy := r.URL.Query().Get("sortBy")

	movies, err := h.MovieService.GetAll(ctx, sortBy)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.GetAllMoviesResponse{
		Data: movies,
	}

	res.ResJson(w, data, http.StatusOK)

}

func (h *MovieHandler) GetMovieByID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	param := r.PathValue("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		res.ErrResJson(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	movie, err := h.MovieService.GetByID(ctx, uint(id))
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res.ResJson(w, movie, http.StatusOK)

}

func (h *MovieHandler) FullUpdateMovieByID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	param := r.PathValue("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		res.ErrResJson(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	body, err := req.DecodedAndValidatedBody[payload.MoviePayload](r.Body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.MovieService.FullUpdate(ctx, uint(id), &body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.MovieResponse{
		Message: "Movie was updated",
	}

	res.ResJson(w, data, http.StatusOK)

}

func (h *MovieHandler) PartialUpdateMovieByID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	param := r.PathValue("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		res.ErrResJson(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	body, err := req.DecodedAndValidatedBody[payload.UpdatePartialMoviePayload](r.Body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.MovieService.PartialUpdate(ctx, uint(id), &body)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.MovieResponse{
		Message: "Movie was updated",
	}

	res.ResJson(w, data, http.StatusOK)
}

func (h *MovieHandler) DeleteMovieByID(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	param := r.PathValue("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		res.ErrResJson(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	err = h.MovieService.Delete(ctx, uint(id))
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.MovieResponse{
		Message: "Movie was deleted",
	}

	res.ResJson(w, data, http.StatusOK)

}

func (h *MovieHandler) SearchMovieByTitle(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	title := r.URL.Query().Get("title")

	movies, err := h.MovieService.SearchMovieByTitle(ctx, title)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.GetAllMoviesResponse{
		Data: movies,
	}

	res.ResJson(w, data, http.StatusOK)
}

func (h *MovieHandler) SearchMovieByActorName(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	actorName := r.URL.Query().Get("actorName")

	movies, err := h.MovieService.SearchMovieByActorName(ctx, actorName)
	if err != nil {
		res.ErrResJson(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &payload.GetAllMoviesResponse{
		Data: movies,
	}

	res.ResJson(w, data, http.StatusOK)
}
