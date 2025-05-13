package service

import (
	"actors-service/internal/model"
	"actors-service/internal/payload"
	"actors-service/internal/repository"
	"context"
)

type ActorService struct {
	ActorRepository *repository.ActorRepository
}

func NewActorService(actorRepository *repository.ActorRepository) *ActorService {
	return &ActorService{
		ActorRepository: actorRepository,
	}
}

func (s *ActorService) Create(ctx context.Context, p *payload.ActorPayload) (uint, error) {
	actorID, err := s.ActorRepository.Create(ctx, p)
	if err != nil {
		return 0, err
	}
	return actorID, nil
}

func (s *ActorService) GetActorsWithMovies(ctx context.Context) ([]model.ActorWithMovies, error) {
	actors, err := s.ActorRepository.GetActorsWithMovies(ctx)
	if err != nil {
		return nil, err
	}
	return actors, nil
}

func (s *ActorService) GetByID(ctx context.Context, id uint) (*model.Actor, error) {
	actor, err := s.ActorRepository.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return actor, nil
}

func (s *ActorService) FullUpdate(ctx context.Context, id uint, p *payload.ActorPayload) error {
	err := s.ActorRepository.FullUpdate(ctx, id, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *ActorService) PartialUpdate(ctx context.Context, id uint, p *payload.PartialUpdateActorPayload) error {
	err := s.ActorRepository.PartialUpdate(ctx, id, p)
	if err != nil {
		return err
	}
	return nil
}

func (s *ActorService) Delete(ctx context.Context, id uint) error {
	err := s.ActorRepository.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
