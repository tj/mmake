package installer_test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/tj/mmake/installer"
	"github.com/tj/mmake/resolver"
)

func init() {
	log.SetLevel(log.WarnLevel)
	log.SetHandler(cli.Default)
}

// Tree shim.
func tree(dir string) (string, error) {
	cmd := exec.Command("tree", dir)
	b, err := cmd.CombinedOutput()
	return string(b), err
}

// Remove junk.
func remove() error {
	return os.RemoveAll("/tmp/include/github.com/tj")
}

func ExampleInstaller_Parse() {
	remove()

	config := installer.Config{
		Resolver:    resolver.NewGithubResolver(),
		Destination: "/tmp/include",
		Log:         log.Log,
	}

	i := installer.New(config)

	r := strings.NewReader(`
include github.com/tj/foo/foo
include github.com/tj/foo/bar
include github.com/tj/foo/foo.mk
include something.mk
	`)

	err := i.Parse(r)

	if err != nil {
		log.WithError(err).Fatal("installing")
	}

	s, err := tree("/tmp/include/github.com")
	if err != nil {
		log.WithError(err).Fatal("tree")
	}

	fmt.Printf("%s\n", s)

	// Output:
	// github.com
	// └── tj
	//     └── foo
	//         ├── bar
	//         ├── foo
	//         └── foo.mk
	//
	// 3 directories, 3 files
}

func ExampleInstaller_Install() {
	remove()

	config := installer.Config{
		Resolver:    resolver.NewGithubResolver(),
		Destination: "/tmp/include",
		Log:         log.Log,
	}

	i := installer.New(config)

	err := i.Install([]string{
		"github.com/tj/foo/foo.mk",
		"github.com/tj/foo/bar",
	})

	if err != nil {
		log.WithError(err).Fatal("installing")
	}

	s, err := tree("/tmp/include/github.com")
	if err != nil {
		log.WithError(err).Fatal("tree")
	}

	fmt.Printf("%s\n", s)

	// Output:
	// github.com
	// └── tj
	//     └── foo
	//         ├── bar
	//         └── foo.mk
	//
	// 3 directories, 2 files
}

func ExampleInstaller_Install_many() {
	remove()

	config := installer.Config{
		Resolver:    resolver.NewGithubResolver(),
		Destination: "/tmp/include",
		Log:         log.Log,
	}

	i := installer.New(config)

	err := i.Install([]string{
		"github.com/tj/foo/bar",
		"github.com/tj/foo/foo.mk",
		"github.com/tj/foo/some/nested/stuff.mk",
		"github.com/tj/foo/deps/foo",
	})

	if err != nil {
		log.WithError(err).Fatal("installing")
	}

	s, err := tree("/tmp/include/github.com")
	if err != nil {
		log.WithError(err).Fatal("tree")
	}

	fmt.Printf("%s\n", s)

	// Output:
	// github.com
	// └── tj
	//     └── foo
	//         ├── bar
	//         ├── deps
	//         │   ├── bar.mk
	//         │   ├── baz
	//         │   └── foo
	//         ├── foo.mk
	//         └── some
	//             └── nested
	//                 └── stuff.mk
	//
	// 6 directories, 6 files
}
