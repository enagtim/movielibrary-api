package repository

import (
	"context"
	"database/sql"
	"movies-service/internal/model"
	"movies-service/internal/payload"
	"movies-service/internal/postgres"
	"movies-service/pkg/consts"

	sq "github.com/Masterminds/squirrel"
)

type MovieRepository struct {
	Database *postgres.Db
}

func NewMovieRepository(db *postgres.Db) *MovieRepository {
	return &MovieRepository{Database: db}
}

func (r *MovieRepository) Create(ctx context.Context, p *payload.MoviePayload) (uint, error) {
	query, args, err := sq.
		Insert("movies").
		Columns("title", "description", "release_date", "rating").
		Values(p.Title, p.Description, p.ReleaseDate, p.Rating).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, consts.ErrFailedToBuildSQL
	}
	var movieID uint

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&movieID)

	if err != nil {
		return 0, consts.ErrFailedCreateMovie
	}

	return movieID, nil

}

func (r *MovieRepository) GetByID(ctx context.Context, id uint) (*model.Movie, error) {
	var movie model.Movie

	query, args, err := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.Title, &movie.Decription, &movie.ReleaseDate, &movie.Rating)

	if err == sql.ErrNoRows {
		return nil, consts.ErrMovieNotFound
	}

	return &movie, nil
}

func (r *MovieRepository) FullUpdate(ctx context.Context, id uint, p *payload.MoviePayload) error {
	query, args, err := sq.
		Update("movies").
		Set("title", p.Title).
		Set("description", p.Description).
		Set("release_date", p.ReleaseDate).
		Set("rating", p.Rating).
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return consts.ErrFailedToBuildSQL
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return consts.ErrFailedUpdateMovie
	}

	return nil
}

func (r *MovieRepository) PartialUpdate(ctx context.Context, id uint, p *payload.UpdatePartialMoviePayload) error {
	updateBuilder := sq.Update("movies").Where(sq.Eq{"id": id})

	if p.Title != nil {
		updateBuilder = updateBuilder.Set("title", *p.Title)
	}
	if p.Description != nil {
		updateBuilder = updateBuilder.Set("description", *p.Description)
	}

	if p.ReleaseDate != nil {
		updateBuilder = updateBuilder.Set("release_date", *p.ReleaseDate)
	}

	if p.Rating != nil {
		updateBuilder = updateBuilder.Set("rating", *p.Rating)
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

func (r *MovieRepository) Delete(ctx context.Context, id uint) error {

	query, args, err := sq.
		Delete("movies").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return consts.ErrFailedToBuildSQL
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return consts.ErrFailedDeleteMovie
	}

	return nil

}
