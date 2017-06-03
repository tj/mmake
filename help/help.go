// Package help provides target help output.
package help

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/tj/mmake/parser"
)

// OutputAllShort outputs all short help representations to the given writer.
func OutputAllShort(r io.Reader, w io.Writer, targets []string) error {
	comments, err := getComments(r, targets)
	if err != nil {
		return err
	}

	width := targetWidth(comments)
	fmt.Fprintf(w, "\n")
	for _, c := range comments {
		if c.Target == "" {
			continue
		}

		fmt.Fprintf(w, "  %-*s %-s\n", width+2, c.Target, firstLine(c.Value))
	}

	fmt.Fprintf(w, "\n")
	return nil
}

// OutputAllLong outputs all long help representations to the given writer.
func OutputAllLong(r io.Reader, w io.Writer, targets []string) error {
	comments, err := getComments(r, targets)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "\n")
	for _, c := range comments {
		if c.Target == "" {
			continue
		}

		fmt.Fprintf(w, "  %-s:\n%-s\n\n", c.Target, indent(indent(c.Value)))
	}

	fmt.Fprintf(w, "\n")
	return nil
}

// getComments parses, filters, and sorts all comment nodes.
func getComments(r io.Reader, targets []string) ([]parser.Comment, error) {
	nodes, err := parser.ParseRecursive(r, "/usr/local/include")
	if err != nil {
		return nil, errors.Wrap(err, "parsing")
	}

	comments := filterComments(nodes, targets)
	sort.Sort(byTarget(comments))
	return comments, nil
}

// Filter comment nodes.
func filterComments(nodes []parser.Node, targets []string) (comments []parser.Comment) {
Outer:
	for _, n := range nodes {
		c, ok := n.(parser.Comment)

		if !ok {
			continue
		}

		if len(targets) == 0 {
			comments = append(comments, c)
			continue
		}

		for _, t := range targets {
			if match, _ := filepath.Match(t, c.Target); match {
				comments = append(comments, c)
				continue Outer
			}
		}
	}

	return
}

// Sort by comment target string.
type byTarget []parser.Comment

func (v byTarget) Len() int           { return len(v) }
func (v byTarget) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byTarget) Less(i, j int) bool { return v[i].Target < v[j].Target }

// Target width from the given comments.
func targetWidth(comments []parser.Comment) (n int) {
	for _, c := range comments {
		if len(c.Target) > n {
			n = len(c.Target)
		}
	}

	return
}

// First line of `s`.
func firstLine(s string) string {
	return strings.Split(s, "\n")[0]
}

// Indent the given string.
func indent(s string) string {
	return strings.Replace("  "+s, "\n", "\n  ", -1)
}
