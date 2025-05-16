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
	tx, err := r.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, consts.ErrFailedToBeginTx
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	query, args, err := sq.
		Insert("movies").
		Columns("title", "description", "release_date", "rating").
		Values(p.Title, p.Description, p.ReleaseDate, p.Rating).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, consts.ErrFailedToBuildSQL
	}
	var movieID uint

	err = tx.QueryRowContext(ctx, query, args...).Scan(&movieID)
	if err != nil {
		return 0, consts.ErrFailedCreateMovie
	}

	for _, actorID := range p.ActorsIDs {
		linkQuery, linkArgs, err := sq.
			Insert("movie_actors").
			Columns("movie_id", "actor_id").
			Values(movieID, actorID).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return 0, consts.ErrFailedToBuildSQL
		}

		_, err = tx.ExecContext(ctx, linkQuery, linkArgs...)
		if err != nil {
			return 0, consts.ErrFailedToLinkActors
		}
	}

	return movieID, nil

}

func (r *MovieRepository) GetAll(ctx context.Context, sortBy string) ([]model.Movie, error) {
	validateSortFields := map[string]string{
		"title":        "title",
		"release_date": "release_date",
		"rating":       "rating",
	}
	sortField, ok := validateSortFields[sortBy]

	if !ok {
		sortField = "rating"
	}

	query, args, err := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		OrderBy(sortField + " DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	rows, err := r.Database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, consts.ErrFailedToExecute
	}

	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var m model.Movie
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Description,
			&m.ReleaseDate,
			&m.Rating)
		if err != nil {
			return nil, consts.ErrFailedToScanRow
		}
		movies = append(movies, m)
	}
	return movies, nil

}

func (r *MovieRepository) GetByID(ctx context.Context, id uint) (*model.Movie, error) {
	var movie model.Movie

	query, args, err := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(
		&movie.ID,
		&movie.Title,
		&movie.Description,
		&movie.ReleaseDate,
		&movie.Rating)
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
		PlaceholderFormat(sq.Dollar).
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

	if p.ActorsIDs != nil {
		updateBuilder = updateBuilder.Set("actors_ids", *p.ActorsIDs)
	}

	query, args, err := updateBuilder.PlaceholderFormat(sq.Dollar).ToSql()
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
		PlaceholderFormat(sq.Dollar).
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

func (r *MovieRepository) SearchMovieByTitle(ctx context.Context, title string) ([]model.Movie, error) {
	query, args, err := sq.
		Select("id", "title", "description", "release_date", "rating").
		From("movies").
		Where(sq.Like{"title": "%" + title + "%"}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	rows, err := r.Database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, consts.ErrFailedToExecute
	}

	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		if err != nil {
			return nil, consts.ErrFailedToScanRow
		}
		movies = append(movies, movie)
	}

	err = rows.Err()
	if err != nil {
		return nil, consts.ErrFailedToProcessRows
	}

	return movies, nil
}

func (r *MovieRepository) SearchMovieByActorName(ctx context.Context, actorName string) ([]model.Movie, error) {
	query, args, err := sq.
		Select("movies.id", "movies.title", "movies.description", "movies.release_date", "movies.rating").
		From("movies").
		Join("movie_actors ON movie_actors.movie_id = movies.id").
		Join("actors ON actors.id = movie_actors.actor_id").
		Where(sq.Like{"actors.name": "%" + actorName + "%"}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	rows, err := r.Database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, consts.ErrFailedToExecute
	}

	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.Description,
			&movie.ReleaseDate,
			&movie.Rating,
		)
		if err != nil {
			return nil, consts.ErrFailedToScanRow
		}
		movies = append(movies, movie)
	}

	err = rows.Err()
	if err != nil {
		return nil, consts.ErrFailedToProcessRows
	}

	return movies, nil

}
