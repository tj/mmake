package resolver

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// github implementation.
type github struct{}

// NewGithubResolver returns a github resolver.
func NewGithubResolver() Interface {
	return &github{}
}

// Get implementation.
func (r *github) Get(s string) (io.ReadCloser, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s", s))
	if err != nil {
		return nil, errors.Wrap(err, "parsing include path")
	}

	if u.Host != "github.com" {
		return nil, ErrNotSupported
	}

	parts := strings.SplitN(u.Path, "/", 4)
	if len(parts) < 3 {
		return nil, errors.New("user, repo required in include url")
	}

	if len(parts) < 4 {
		parts = append(parts, "index.mk")
	}

	user := parts[1]
	repo := parts[2]
	file := parts[3]
	raw := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", user, repo, file)

	return r.fetch(raw)
}

// Fetch over HTTP.
func (r *github) fetch(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "requesting")
	}

	if res.StatusCode == 200 {
		return res.Body, nil
	}

	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrNotFound
	}

	return nil, errors.Errorf("response %s", res.Status)
}
