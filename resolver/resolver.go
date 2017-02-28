package resolver

import (
	"errors"
	"io"
)

// Errors.
var (
	// ErrNotFound is returned when a non-local lookup fails.
	ErrNotFound = errors.New("not found")

	// ErrNotSupported is returned when an unsupported include is detected,
	// and may be checked by subsequent resolvers, or ignored completely.
	ErrNotSupported = errors.New("unsupported")
)

// Interface for resolving and fetching includes.
type Interface interface {
	// Get must return the contents of the file,
	// ErrNotFound when the file is not found,
	// or an error.
	Get(path string) (io.ReadCloser, error)
}
