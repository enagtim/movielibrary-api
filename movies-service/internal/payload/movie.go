package payload

import (
	"movies-service/internal/model"
	"time"
)

type MoviePayload struct {
	Title       string    `json:"title" validate:"required,min=1,max=150"`
	Description *string   `json:"description" validate:"max=1000"`
	ReleaseDate time.Time `json:"release_date" validate:"required"`
	Rating      float64   `json:"rating" validate:"gte=0,lte=10"`
	ActorsIDs   []uint    `json:"actors_ids" validate:"required"`
}

type UpdatePartialMoviePayload struct {
	Title       *string    `json:"title" validate:"min=1,max=150"`
	Description *string    `json:"description" validate:"max=1000"`
	ReleaseDate *time.Time `json:"release_date"`
	Rating      *float64   `json:"rating" validate:"gte=0,lte=10"`
	ActorsIDs   *[]uint    `json:"actors_ids"`
}

type CreateMovieResponse struct {
	MovieID uint   `json:"movieID"`
	Message string `json:"message"`
}

type MovieResponse struct {
	Message string `json:"message"`
}

type GetAllMoviesResponse struct {
	Data []model.Movie `json:"data"`
}
