package consts

import "errors"

var (
	ErrFailedToBuildSQL    = errors.New("failed to build SQL query")
	ErrFailedCreateActor   = errors.New("failed to create actor")
	ErrActorNotFound       = errors.New("actor not found")
	ErrFailedUpdateActor   = errors.New("failed to update actor")
	ErrFailedToExecute     = errors.New("failed to execute query")
	ErrInvalidAffectedrows = errors.New("failed to get affected rows")
	ErrFailedDeleteActor   = errors.New("failed to delete actor")
)
