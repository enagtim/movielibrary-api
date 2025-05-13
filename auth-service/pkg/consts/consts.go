package consts

import "errors"

var (
	ErrFailedToBuildSQL        = errors.New("failed to build SQL query")
	ErrFailedCreateUser        = errors.New("failed to create user")
	ErrFailedHashedPassword    = errors.New("failed hashed password")
	ErrFailedGetUserByUserName = errors.New("failed get user by username")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrGenerateToken           = errors.New("error generate jwt token")
)
