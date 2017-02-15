![Modern Make](https://dl.dropboxusercontent.com/u/6396913/mmake/gh-title.png)

## About

Mmake is a small program which wraps `make` to provide additional functionality, such as user-friendly help output, remote includes,
and eventually more.

## Installation

Grab a [binary]() or:

```
$ go get github.com/tj/mmake/cmd/mmake
```

Next add the following alias to your profile:

```
alias make=mmake
```

## Features

- Remote includes
- Target help output

### Help output

Make's primary function is not to serve as a "task runner", however it's often used for that scenario due to its ubiquitous nature.

Since it was not designed for this, its support for outputting target (or "task") help documentation" does not really exist.

Suppose you have the following makefile, standard Make has no notion comments tied to a given target, but wrapping make can provide this.

```Makefile
# Start the dev server.
#
# Note that the API server must
# also be running.
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
size:
	@gopherjs build client/*.go -m -o /tmp/out.js
	@du -h /tmp/out.js
	@gopher-count /tmp/out.js | sort -nr
.PHONY: size

```

For example output all target docs:

```
$ alias make=mmake
$ make help

  start      Start the dev server.
  api        Start the API server.
  deps       Display dependency graph.
  size       Display size of dependencies.

```

Or output verbose help output of a single target:

```
$ make help start

  Start the dev server.

  Note that the API server must
  also be running.

```

The default behaviour of Make is of course preserved:

```
$ make
serving at http://localhost:3000 and on port 3000 of any available addresses

$ make size
...
```

### Remote includes

Includes may specify a URL for inclusion, which are automatically downloaded to /usr/local/include and become available to Make. Note that make resolves includes to this directory by default, so the Makefile will still work for regular users.

Includes are resolved recursively. For example you may have a standard set of includes for your team to run tests, lint, and deploy:

```Makefile
include github.com/apex/make/deploy
include github.com/apex/make/lint
include github.com/apex/make/test
```

This can be a lot to remember, so you could also provide a file which includes the others:

```Makefile
include github.com/apex/make/all
```

Or perhaps one per dev environment such as Node or Golang:

```Makefile
include github.com/apex/make/node
include github.com/apex/make/golang
```

If you're worried about arbitrary code execution, then simply fork a project and maintain keep the code private.

## Registry

If you're looking to find or share makefiles check out the [Wiki](https://github.com/tj/mmake/wiki/Registry), and feel free to add a category if it is missing.

## Links

- [Wiki / Registry](https://github.com/tj/mmake/wiki/Registry)

## Badges

[![GoDoc](https://godoc.org/github.com/tj/mmake?status.svg)](https://godoc.org/github.com/tj/mmake)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
[![](http://apex.sh/images/badge.svg)](https://apex.sh/ping/)

---

> [tjholowaychuk.com](http://tjholowaychuk.com) &nbsp;&middot;&nbsp;
> GitHub [@tj](https://github.com/tj) &nbsp;&middot;&nbsp;
> Twitter [@tjholowaychuk](https://twitter.com/tjholowaychuk)
