package repository

import (
	"actors-service/internal/model"
	"actors-service/pkg/db"
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

type ActorRepository struct {
	Database *db.Db
}

func NewActorRepository(db *db.Db) *ActorRepository {
	return &ActorRepository{Database: db}
}

func (r *ActorRepository) Create(ctx context.Context, actor model.Actor) (uint, error) {
	query, args, err := sq.
		Insert("actors").
		Columns("name", "gender", "birth_date").
		Values(actor.Name, actor.Gender, actor.BirthDate).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, errors.New("failed to build query")
	}

	var actorID uint

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&actorID)

	if err != nil {
		return 0, errors.New("failed to create actor")
	}

	return actorID, nil
}

func (r *ActorRepository) GetById(ctx context.Context, id uint) (*model.Actor, error) {
	var actor model.Actor

	query, args, err := sq.
		Select("id", "name", "gender", "birth_date").
		From("actors").
		Where()
}
