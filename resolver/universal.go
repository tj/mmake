package resolver

import (
	"io"
	"strings"
)

// universal implementation.
type universal struct{}

// NewUniversalResolver returns a universal resolver.
func NewUniversalResolver() Interface {
	return &universal{}
}

// Get implementation.
// Chooses an appropriate resolver
func (r *universal) Get(s string) (io.ReadCloser, error) {
	if strings.HasPrefix(s, "http") {
		return NewHTTPResolver().Get(s)
	} else if strings.HasPrefix(s, "github.com") {
		return NewGithubResolver().Get(s)
	} else {
		return NewLocalResolver().Get(s)
	}
}
