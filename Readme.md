# Modern Make

## About

Mmake is a small program which wraps `make` to provide additional functionality, such as user-friendly help output, remote includes,
and eventually more. It otherwise acts as a pass-through to standard make.

## Installation

Mechanisms:

Grab a [binary](https://github.com/tj/mmake/releases)

Build from source:
```
$ go get github.com/tj/mmake/cmd/mmake
```

Homebrew:
```
$ brew tap tj/mmake https://github.com/tj/mmake.git
$ brew install tj/mmake/mmake
```

Next add the following alias to your profile:

```
alias make=mmake
```

## Features

### Help output

Make's primary function is not to serve as a "task runner", however it's often used for that scenario due to its ubiquitous nature, and if you're already using it, why not! Make is however lacking a built-in mechanism for displaying help information.

Here's an example Makefile:

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
#- Any comment preceded by a dash is omitted.
size:
	@gopherjs build client/*.go -m -o /tmp/out.js
	@du -h /tmp/out.js
	@gopher-count /tmp/out.js | sort -nr
.PHONY: size

```

Mmake provides a `help` command to display all target comments in short form:

```
$ alias make=mmake
$ make help

  start      Start the dev server.
  api        Start the API server.
  deps       Display dependency graph.
  size       Display size of dependencies.

```

You can optionally filter which commands to view the help dialogue for (this supports [standard Unix glob patterns](https://en.wikipedia.org/wiki/Glob_(programming)#Syntax)):

```
$ make help start

  start   Start the dev server.

$ make help s*

  size    Display size of dependencies.
  start   Start the dev server.

```

The `help` command also supports displaying longer output with the verbose flag (`-v` / `--verbose`):

```
$ make help -v start

  Start the dev server.

  Note that the API server must
  also be running.

```

```
$ make help -v

  start:
    Start the dev server.

    Note that the API server must
    also be running.

  api:    
    Start the API server.

  deps:       
    Display dependency graph.

  size:
    Display size of dependencies.
    
```

The default behaviour of Make is of course preserved:

```
$ make
serving at http://localhost:3000 and on port 3000 of any available addresses

$ make size
...
```

### Remote includes

Includes may specify a URL (http, https, or github shortcut) for inclusion, which are automatically downloaded to `/usr/local/include` and become available to Make. Note that make resolves includes to this directory by default, so the Makefile will still work for regular users.

Includes are resolved recursively. For example you may have a standard set of includes for your team to run tests, lint, and deploy:

```Makefile
include github.com/apex/make/deploy
include github.com/apex/make/lint
include github.com/apex/make/test
include https://github.com/apex/make/test/Makefile
include https://github.com/apex/make/test/make.mk
```

This can be a lot to remember, so you could also provide a file which includes the others:

```Makefile
include github.com/apex/make/all
```

If the given repository contains an `index.mk` file, you can just declare:

```Makefile
include github.com/apex/make
```

Or perhaps one per dev environment such as Node or Golang:

```Makefile
include github.com/apex/make/node
include github.com/apex/make/golang
```

If you're worried about arbitrary code execution, then simply fork a project and maintain control over it.

#### Update

Once the remote includes are downloaded to `/usr/local/include`, `mmake` will not try to fetch them again. In order to get an updated copy of the remote includes, `mmake` provides an `update` target that will download them again:

```
$ make update
```

## Registry

If you're looking to find or share makefiles check out the [Wiki](https://github.com/tj/mmake/wiki/Registry), and feel free to add a category if it is missing.

## Links

- [GNU Make](https://www.gnu.org/software/make/manual/make.html) documentation
- [Wiki](https://github.com/tj/mmake/wiki/Registry) registry
- [Announcement](https://medium.com/@tjholowaychuk/modern-make-b55d53cf80d9#.q1u1knrf5) blog post
- [Introduction](https://www.youtube.com/watch?v=NLS_gbg4_wI) youtube video
- [AUR Package](https://aur.archlinux.org/packages/mmake-bin/) Arch Linux Package

## Badges

[![GoDoc](https://godoc.org/github.com/tj/mmake?status.svg)](https://godoc.org/github.com/tj/mmake)
![](https://img.shields.io/badge/license-MIT-blue.svg)
![](https://img.shields.io/badge/status-stable-green.svg)
[![](http://apex.sh/images/badge.svg)](https://apex.sh/ping/)

---

> [tjholowaychuk.com](http://tjholowaychuk.com) &nbsp;&middot;&nbsp;
> GitHub [@tj](https://github.com/tj) &nbsp;&middot;&nbsp;
> Twitter [@tjholowaychuk](https://twitter.com/tjholowaychuk)
