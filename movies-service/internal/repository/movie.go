package repository

import (
	"context"
	"database/sql"
	"errors"
	"movies-service/internal/model"
	"movies-service/internal/payload"
	"movies-service/internal/postgres"

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
		return 0, errors.New("failed to build query")
	}
	var movieID uint

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&movieID)

	if err != nil {
		return 0, errors.New("failed to create movie")
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
		return nil, errors.New("failed to build query")
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.Title, &movie.Decription, &movie.ReleaseDate, &movie.Rating)

	if err == sql.ErrNoRows {
		return nil, errors.New("failed to get movie")
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
		return errors.New("failed to build query")
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return errors.New("failed to update movie")
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
		return errors.New("failed to build query")
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

func (r *MovieRepository) Delete(ctx context.Context, id uint) error {

	query, args, err := sq.
		Delete("movies").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		return errors.New("failed to build query")
	}

	_, err = r.Database.DB.ExecContext(ctx, query, args...)

	if err != nil {
		return errors.New("failed to delete movie")
	}

	return nil

}
