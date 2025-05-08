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

func (r *ActorRepository) Create(ctx context.Context, p *payload.ActorPaylod) (uint, error) {
	query, args, err := sq.
		Insert("actors").
		Columns("name", "gender", "birth_date").
		Values(p.Name, p.Gender, p.BirthDate).
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

func (r *ActorRepository) FullUpdate(ctx context.Context, id uint, p *payload.ActorPaylod) error {
	query, args, err := sq.
		Update("actors").Set("name", p.Name).
		Set("gender", p.Gender).
		Set("birth_date", p.BirthDate).
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

func (r *ActorRepository) PartialhUpdate(ctx context.Context, id uint, p *payload.PartialUpdateActorPaylod) error {
	updateBuilder := sq.Update("actors").Where(sq.Eq{"id": id})

	if p.Name != nil {
		updateBuilder = updateBuilder.Set("name", *p.Name)
	}
	if p.Gender != nil {
		updateBuilder = updateBuilder.Set("gender", *p.Gender)
	}
	if p.BirthDate != nil {
		updateBuilder = updateBuilder.Set("birth_date", *p.BirthDate)
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
