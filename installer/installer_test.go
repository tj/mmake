package installer_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

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

type mockResolver struct {
	getFn func(string) (io.ReadCloser, error)
}

func (m mockResolver) Get(s string) (io.ReadCloser, error) { return m.getFn(s) }

func TestInstaller_Update(t *testing.T) {
	remove()

	var resolved bool
	m := mockResolver{
		getFn: func(path string) (io.ReadCloser, error) {
			resolved = true
			return resolver.NewGithubResolver().Get(path)
		},
	}

	config := installer.Config{
		Resolver:    m,
		Destination: "/tmp/include",
		Log:         log.Log,
	}

	i := installer.New(config)

	checkResolve := func(expect bool, update bool, msg string) {
		i.Config.ForceUpdate = update

		err := i.Install([]string{"github.com/tj/foo/bar"})
		if err != nil {
			log.WithError(err).Fatal("installing")
		}

		if resolved != expect {
			t.Fatal(msg)
		}
		resolved = false
	}

	// Doesn't exist, should resolve
	checkResolve(true, false, "Dependency should resolve when it doesn't exist yet")

	// Exists, shouldn't resolve
	checkResolve(false, false, "Dependency should not resolve when it already exists")

	// Exists but update = true, should resolve
	checkResolve(true, true, "Dependency should always resolve when update = true")
}
