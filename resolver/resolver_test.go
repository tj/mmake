package resolver_test

import (
	"io/ioutil"
	"testing"

	"github.com/tj/mmake/resolver"
)

func TestGetIncludePath(t *testing.T) {
	var cases = []struct {
		Args     []string
		Expected string
	}{
		{[]string{}, resolver.DefaultIncludePath},
		{[]string{"update", "-I", "./relative/path"}, "./relative/path"},
		{[]string{"update", "-I./other/path"}, "./other/path"},
		{[]string{"update", "-I./other/path", "-I", "multiple/"}, "./other/path"},
		{[]string{"-I", "multiple/"}, "multiple/"},
		{[]string{"-I"}, resolver.DefaultIncludePath},
		{[]string{"update"}, resolver.DefaultIncludePath},
	}

	for _, c := range cases {
		t.Run(c.Expected, func(t *testing.T) {
			out := resolver.GetIncludePath(c.Args)
			if out != c.Expected {
				t.Errorf("expected %q, got %q", c.Expected, out)
			}
		})
	}
}

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

func TestLocal(t *testing.T) {
	resolver := resolver.NewLocalResolver()

	var cases = []struct {
		Path    string
		Content string
	}{
		{"fixtures/bar", "$(info bar)\n"},
		{"fixtures/foo.mk", "$(info foo)\n"},
		{"fixtures/baz/stuff.mk", "$(info baz)\n"},
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

func TestUniversal(t *testing.T) {
	resolver := resolver.NewUniversalResolver()

	var cases = []struct {
		Path    string
		Content string
	}{
		{"fixtures/bar", "$(info bar)\n"},
		{"github.com/tj/foo/bar", "$(info bar)\n"},
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
