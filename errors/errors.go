package errors

import "errors"

var (
	ErrNotFound = errors.New("document not found")
	ErrNoID     = errors.New("model has no valid ID")
)

// IsNotFound checks if an error is ErrNotFound
func IsNotFound(err error) bool {
	return err == ErrNotFound
}

// IsNoID checks if an error is ErrNoID
func IsNoID(err error) bool {
	return err == ErrNoID
}
