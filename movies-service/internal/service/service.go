package service

import (
	"context"
	"movies-service/internal/model"
	"movies-service/internal/payload"
	"movies-service/internal/repository"
)

type MovieService struct {
	MovieRepository *repository.MovieRepository
}

func NewMovieService(movieRepository *repository.MovieRepository) *MovieService {
	return &MovieService{
		MovieRepository: movieRepository,
	}
}

func (s *MovieService) Create(ctx context.Context, p *payload.MoviePayload) (uint, error) {
	movieID, err := s.MovieRepository.Create(ctx, p)
	if err != nil {
		return 0, err
	}
	return movieID, nil
}

func (s *MovieService) GetAll(ctx context.Context, sortBy string) ([]model.Movie, error) {
	movies, err := s.MovieRepository.GetAll(ctx, sortBy)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) GetByID(ctx context.Context, id uint) (*model.Movie, error) {
	movie, err := s.MovieRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) FullUpdate(ctx context.Context, id uint, p *payload.MoviePayload) error {
	err := s.MovieRepository.FullUpdate(ctx, id, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *MovieService) PartialUpdate(ctx context.Context, id uint, p *payload.UpdatePartialMoviePayload) error {
	err := s.MovieRepository.PartialUpdate(ctx, id, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *MovieService) Delete(ctx context.Context, id uint) error {
	err := s.MovieRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *MovieService) SearchMovieByTitle(ctx context.Context, title string) ([]model.Movie, error) {
	movies, err := s.MovieRepository.SearchMovieByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *MovieService) SearchMovieByActorName(ctx context.Context, actorName string) ([]model.Movie, error) {
	movies, err := s.MovieRepository.SearchMovieByActorName(ctx, actorName)
	if err != nil {
		return nil, err
	}
	return movies, nil
}
