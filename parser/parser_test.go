package parser_test

import (
	"fmt"
	"strings"

	"github.com/tj/mmake/parser"
)

func ExampleParser_Parse_withComments() {
	contents := `
include github.com/tj/foo

# Stuff here:
#
#    :)
#

# Start the dev server.
start:
	@gopherjs -m -v serve --http :3000 github.com/tj/docs/client
.PHONY: start

# Start the API server.
api:
	@go run server/cmd/api/api.go
.PHONY: api

# Display dependency graph.
deps:
	@godepgraph github.com/tj/docs/client | dot -Tsvg | browser
.PHONY: deps

# Display size of dependencies.
#
# - foo
# - bar
# - baz
#
size:
	@gopherjs build client/*.go -m -o /tmp/out.js
	@du -h /tmp/out.js
	@gopher-count /tmp/out.js | sort -nr
.PHONY: size

`

	p := parser.New()

	nodes, err := p.Parse(strings.NewReader(contents))
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		fmt.Printf("%#v\n", node)
	}

	// Output:
	// parser.Include{Value:"github.com/tj/foo"}
	// parser.Comment{Target:"", Value:"Stuff here:\n\n   :)"}
	// parser.Comment{Target:"start", Value:"Start the dev server."}
	// parser.Comment{Target:"api", Value:"Start the API server."}
	// parser.Comment{Target:"deps", Value:"Display dependency graph."}
	// parser.Comment{Target:"size", Value:"Display size of dependencies.\n\n- foo\n- bar\n- baz"}
}

func ExampleParser_Parse_withoutComments() {
	contents := `
include github.com/tj/foo
include github.com/tj/bar

include github.com/tj/something/here

start:
	@gopherjs -m -v serve --http :3000 github.com/tj/docs/client
.PHONY: start

api:
	@go run server/cmd/api/api.go
.PHONY: api

deps:
	@godepgraph github.com/tj/docs/client | dot -Tsvg | browser
.PHONY: deps
`

	p := parser.New()

	nodes, err := p.Parse(strings.NewReader(contents))
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		fmt.Printf("%#v\n", node)
	}

	// Output:
	// parser.Include{Value:"github.com/tj/foo"}
	// parser.Include{Value:"github.com/tj/bar"}
	// parser.Include{Value:"github.com/tj/something/here"}
}
