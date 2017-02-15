# Run all tests.
test:
	@go test -cover ./...
.PHONY: test

# Install the program.
#
# By default the program is installed to
# the $GOPATH/bin directory, thus you must
# have it within your $PATH.
install:
	@go install ./...
.PHONY: install
