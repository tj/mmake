package parser

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Parse the given input.
func Parse(r io.Reader) ([]Node, error) {
	return (&Parser{}).Parse(r)
}

// ParseRecursive parses the given input recursively
// relative to the given dir such as /usr/local/include.
func ParseRecursive(r io.Reader, dir string) ([]Node, error) {
	nodes, err := Parse(r)
	if err != nil {
		return nil, errors.Wrap(err, "parsing")
	}

	for _, n := range nodes {
		inc, ok := n.(Include)

		if !ok {
			continue
		}

		path := filepath.Join(dir, inc.Value)
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "opening %q", path)
		}

		more, err := ParseRecursive(f, dir)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %q", path)
		}

		nodes = append(nodes, more...)

		f.Close()
	}

	return nodes, nil
}
