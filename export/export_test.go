package export_test

import (
	"fmt"
	"testing"

	"github.com/tj/mmake/export"
)

func TestExport_IncludedMakefilePaths(t *testing.T) {
	files := export.IncludedMakefilePaths("/Users/zph/src/mmake/Makefile")
	expected := []string{
		"/usr/local/include/github.com/tj/make/golang",
		"/usr/local/include/github.com/tj/make/cloc",
		"/usr/local/include/github.com/tj/make/todo",
	}

	if fmt.Sprintf("%+v", files) != fmt.Sprintf("%+v", expected) {
		t.Errorf("Expected include paths did not match actual include paths, got: %+v, want: %+v.", files, expected)
	}
}
