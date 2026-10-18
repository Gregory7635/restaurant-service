package restaurant

import "errors"

var (
	ErrNotFound        = errors.New("restaurant not found")
	ErrInvalidInput    = errors.New("invalid input: name and address are required")
	ErrInvalidRating   = errors.New("invalid rating: must be between 0 and 5")
	ErrCmdRepository   = errors.New("unable to command repository")
	ErrQueryRepository = errors.New("unable to query repository")
)
