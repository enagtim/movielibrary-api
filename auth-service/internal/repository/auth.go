package repository

import (
	"auth-service/internal/model"
	"auth-service/internal/payload"
	"auth-service/internal/postgres"
	"auth-service/pkg/consts"
	"context"

	sq "github.com/Masterminds/squirrel"
)

type AuthRepository struct {
	Database *postgres.Db
}

func NewAuthRepository(db *postgres.Db) *AuthRepository {
	return &AuthRepository{
		Database: db,
	}
}

func (r *AuthRepository) Register(ctx context.Context, p *payload.AuthRegisterPayload) (uint, error) {
	query, args, err := sq.
		Insert("users").
		Columns("username", "password_hash", "role").
		Values(p.Username, p.Password, p.Role).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, consts.ErrFailedToBuildSQL
	}

	var userID uint

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		return 0, consts.ErrFailedCreateUser
	}

	return userID, nil
}

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query, args, err := sq.
		Select("id", "username", "password_hash", "role").
		From("users").
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, consts.ErrFailedToBuildSQL
	}

	err = r.Database.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.UserName,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		return nil, consts.ErrFailedGetUserByUserName
	}

	return &user, nil

}
