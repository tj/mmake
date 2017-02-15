# Run all tests.
test:
	@go test -cover ./...
.PHONY: test

# Install the program.
install:
	@go install ./...
.PHONY: install

# Build release.
build:
	@gox -os="linux darwin windows openbsd" ./...
.PHONY: build
