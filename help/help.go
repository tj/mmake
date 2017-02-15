// Package help provides target help output.
package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"

	"github.com/tj/mmake/parser"
)

// OutputAllShort outputs all short help representations to the given writer.
func OutputAllShort(r io.Reader, w io.Writer) error {
	nodes, err := parser.ParseRecursive(r, "/usr/local/include")
	if err != nil {
		return errors.Wrap(err, "parsing")
	}

	fmt.Fprintf(w, "\n")

	for _, n := range nodes {
		c, ok := n.(parser.Comment)

		if !ok || c.Target == "" {
			continue
		}

		fmt.Fprintf(w, "  %-15s %-s\n", c.Target, firstLine(c.Value))
	}

	fmt.Fprintf(w, "\n")
	return nil
}

// OutputTargetLong outputs long help representation of the given target.
func OutputTargetLong(r io.Reader, w io.Writer, target string) error {
	nodes, err := parser.ParseRecursive(r, "/usr/local/include")
	if err != nil {
		return errors.Wrap(err, "parsing")
	}

	fmt.Fprintf(w, "\n")

	for _, n := range nodes {
		c, ok := n.(parser.Comment)

		if !ok || c.Target != target {
			continue
		}

		fmt.Fprintf(w, "%s\n", indent(c.Value))
		break
	}

	fmt.Fprintf(w, "\n")
	return nil
}

// First line of `s`.
func firstLine(s string) string {
	return strings.Split(s, "\n")[0]
}

// Indent the given string.
func indent(s string) string {
	return strings.Replace("  "+s, "\n", "\n  ", -1)
}
