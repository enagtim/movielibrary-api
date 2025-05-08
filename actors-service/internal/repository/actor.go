package repository

import (
	"actors-service/internal/model"
	"actors-service/internal/payload"
	"actors-service/internal/postgres"
	"actors-service/pkg/consts"
	"context"
	"database/sql"

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
		return 0, consts.ErrFailedToBuildSQL
	}

	var actorID uint

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&actorID)

	if err != nil {
		return 0, consts.ErrFailedCreateActor
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
		return nil, consts.ErrFailedToBuildSQL
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.BirthDate)

	if err == sql.ErrNoRows {
		return nil, consts.ErrActorNotFound
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
		return consts.ErrFailedToBuildSQL
	}
	_, err = r.Database.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return consts.ErrFailedUpdateActor
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
		return consts.ErrFailedToBuildSQL
	}

	res, err := r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return consts.ErrFailedToExecute
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return consts.ErrInvalidAffectedrows
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
		return consts.ErrFailedToBuildSQL
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return consts.ErrFailedDeleteActor
	}

	return nil
}
