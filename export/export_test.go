package export_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tj/mmake/export"
)

func TestExport_IncludedMakefilePaths(t *testing.T) {
	files := export.IncludedMakefilePaths("/Users/zph/src/mmake/Makefile")
	expected := []string{
		"/Users/zph/src/mmake/Makefile",
		"/usr/local/include/github.com/tj/make/golang",
		"/usr/local/include/github.com/tj/make/cloc",
		"/usr/local/include/github.com/tj/make/todo",
	}

	if fmt.Sprintf("%+v", files) != fmt.Sprintf("%+v", expected) {
		t.Errorf("Expected include paths did not match actual include paths, got: %+v, want: %+v.", files, expected)
	}
}

func TestExport_FullRun(t *testing.T) {
	output := `
#- start=include github.com/tj/make/golang
#- start=include github.com/tj/make/cloc
# Output source statistics.
cloc:
        @cloc --exclude-dir=client,vendor .
.PHONY: cloc
#- end=include github.com/tj/make/cloc

#- start=include github.com/tj/make/todo
# Output to-do items per file.
todo:
        @grep \
                --exclude-dir=./vendor \
                --exclude-dir=./client/node_modules \
                --text \
                --color \
                -nRo ' TODO:.*' .
.PHONY: todo
#- end=include github.com/tj/make/todo

# Run all tests.
test:
        @go test -cover ./...
.PHONY: test

# Install the commands.
install:
        @go install ./cmd/...
.PHONY: install

# Release binaries to GitHub.
release:
        @goreleaser --rm-dist --config .goreleaser.yml
.PHONY: release

# Show size of imports.
size:
        @curl -sL https://gist.githubusercontent.com/tj/04e0965e23da00ca33f101e5b2ed4ed4/raw/9aa16698b2bc606cf911219ea540972edef05c4b/gistfile1.txt | bash
.PHONY: size
#- end=include github.com/tj/make/golang
	`
	expected := strings.Split(output, "\n")
	path := "/Users/zph/src/mmake/Makefile"
	includes := export.IncludedMakefilePaths(path)
	actual := export.ProcessMakefile(path, includes, []string{})
	if fmt.Sprintf("%+v", actual) != fmt.Sprintf("%+v", expected) {
		t.Errorf("Expected include paths did not match actual include paths, got: %+v, want: %+v.", actual, expected)
	}

}
