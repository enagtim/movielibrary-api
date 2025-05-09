package consts

import "errors"

var (
	ErrFailedToBuildSQL    = errors.New("failed to build SQL query")
	ErrFailedCreateMovie   = errors.New("failed to create movie")
	ErrMovieNotFound       = errors.New("movie not found")
	ErrFailedUpdateMovie   = errors.New("failed to update movie")
	ErrFailedToExecute     = errors.New("failed to execute query")
	ErrInvalidAffectedrows = errors.New("failed to get affected rows")
	ErrFailedDeleteMovie   = errors.New("failed to delete movie")
	ErrFailedToScanRow     = errors.New("failed to scan row")
	ErrFailedToProcessRows = errors.New("failed to process rows")
	ErrFailedToBeginTx     = errors.New("failed to begin transaction")
	ErrFailedToLinkActors  = errors.New("failed to link actors")
)
