package resolver_test

import (
	"io/ioutil"
	"testing"

	"github.com/tj/mmake/resolver"
)

func TestGithub(t *testing.T) {
	resolver := resolver.NewGithubResolver()

	var cases = []struct {
		Path    string
		Content string
	}{
		{"github.com/tj/foo/bar", "$(info bar)\n"},
		{"github.com/tj/foo/foo.mk", "$(info foo)\n"},
		{"github.com/tj/foo/baz/stuff.mk", "$(info baz)\n"},
	}

	for _, c := range cases {
		t.Run(c.Path, func(t *testing.T) {
			r, err := resolver.Get(c.Path)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			b, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if string(b) != c.Content {
				t.Errorf("expected %q, got %q", c.Content, string(b))
			}
		})
	}
}
