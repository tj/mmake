package resolver

import (
	"errors"
	"io"
)

// DefaultIncludePath is the default directory where make includes reside and
// to which mmake will download remote includes.
const DefaultIncludePath = "/usr/local/include"

// Errors.
var (
	// ErrNotFound is returned when a non-local lookup fails.
	ErrNotFound = errors.New("not found")

	// ErrNotSupported is returned when an unsupported include is detected,
	// and may be checked by subsequent resolvers, or ignored completely.
	ErrNotSupported = errors.New("unsupported")
)

// IncludePath is the directory where make includes reside and to which mmake
// will download remote includes for the duration of a particular execution
// of mmake.
var IncludePath = DefaultIncludePath

// Interface for resolving and fetching includes.
type Interface interface {
	// Get must return the contents of the file,
	// ErrNotFound when the file is not found,
	// or an error.
	Get(path string) (io.ReadCloser, error)
}

// GetIncludePath allows the include path to be overridden.
func GetIncludePath(args []string) string {
	// always reset the global IncludePath to the default
	IncludePath = DefaultIncludePath

	for idx, arg := range args {
		if arg == "-I" && len(args) > idx+1 {
			IncludePath = args[idx+1]
			break
		} else if arg[:2] == "-I" && len(arg) > 2 {
			IncludePath = arg[2:]
			break
		}
	}

	return IncludePath
}
