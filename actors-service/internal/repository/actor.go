package repository

import (
	"actors-service/internal/model"
	"actors-service/internal/payload"
	"actors-service/internal/postgres"
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

type ActorRepository struct {
	Database *postgres.Db
}

func NewActorRepository(db *postgres.Db) *ActorRepository {
	return &ActorRepository{Database: db}
}

func (r *ActorRepository) Create(ctx context.Context, payload payload.ActorPaylod) (uint, error) {
	query, args, err := sq.
		Insert("actors").
		Columns("name", "gender", "birth_date").
		Values(payload.Name, payload.Gender, payload.BirthDate).
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
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, errors.New("failed to build query")
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)

	if err == sql.ErrNoRows {
		return nil, errors.New("actor not found")
	}

	if err != nil {
		return nil, errors.New("failed to get actor")
	}

	return &actor, nil
}

func (r *ActorRepository) PutUpdate(ctx context.Context, id uint, payload payload.ActorPaylod) error {
	query, args, err := sq.
		Update("actors").Set("name", payload.Name).
		Set("gender", payload.Gender).
		Set("birth_date", payload.BirthDate).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return errors.New("failed to build update query")
	}
	_, err = r.Database.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.New("failed to update actor")
	}
	return nil
}

func (r *ActorRepository) PatchUpdate(ctx context.Context, id uint, payload payload.PartialUpdateActorPaylod) error {
	updateBuilder := sq.Update("actors").Where(sq.Eq{"id": id})

	if payload.Name != nil {
		updateBuilder = updateBuilder.Set("name", *payload.Name)
	}
	if payload.Gender != nil {
		updateBuilder = updateBuilder.Set("gender", *payload.Gender)
	}
	if payload.BirthDate != nil {
		updateBuilder = updateBuilder.Set("birth_date", *payload.BirthDate)
	}

	query, args, err := updateBuilder.ToSql()

	if err != nil {
		return errors.New("failed to build update query")
	}

	res, err := r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return errors.New("failed to execute partial update")
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return errors.New("failed to get affected rows")
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil

}

func (r *ActorRepository) Delete(ctx context.Context, id uint) error {
	query, args, err := sq.
		Delete("actors").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return errors.New("failed to build delete query")
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return errors.New("failed to delete actor")
	}

	return nil
}
