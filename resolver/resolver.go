package resolver

import (
	"errors"
	"io"
)

// Errors.
var (
	ErrNotFound = errors.New("not found")
)

// Interface for resolving and fetching includes.
type Interface interface {
	// Get must return the contents of the file,
	// ErrNotFound when the file is not found,
	// or an error.
	Get(path string) (io.ReadCloser, error)
}
