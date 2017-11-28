package resolver

import (
	"io"
	"os"
)

// github implementation.
type local struct{}

// NewGithubResolver returns a github resolver.
func NewLocalResolver() Interface {
	return &local{}
}

// Get implementation.
func (r *local) Get(s string) (io.ReadCloser, error) {
	return os.Open(s)
}
