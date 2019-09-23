package resolver

import (
	"io"
	"net/http"
	"net/url"

	log "github.com/apex/log"
	"github.com/pkg/errors"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

// github implementation.
type httpResolver struct{}

// NewHTTPResolver returns a http resolver.
func NewHTTPResolver() Interface {
	return &httpResolver{}
}

// Get implementation.
func (r *httpResolver) Get(s string) (io.ReadCloser, error) {
	log.WithField("resolver", "http").Debug("Resolver type")
	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.Wrap(err, "parsing include path")
	}

	log.WithField("uri", u).Debug("URI")
	if u.Scheme == "http" || u.Scheme == "https" {
		return r.fetch(s)
	}
	return nil, ErrNotSupported
}

// Fetch over HTTP.
func (r *httpResolver) fetch(url string) (io.ReadCloser, error) {
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
