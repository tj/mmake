package installer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/pkg/errors"

	"github.com/tj/mmake/parser"
	"github.com/tj/mmake/resolver"
)

// TODO: dedupe in-flight
// TODO: tests for errors

// Config for installer.
type Config struct {
	Destination string             // Destination of stored files.
	Resolver    resolver.Interface // Resolver implementation.
	Log         log.Interface      // Log implementation.
}

// Installer fetches and resolves dependencies recursively
// to the configured Destination using the Resolver
// implementation provided.
type Installer struct {
	Config
}

// New installer.
func New(config Config) *Installer {
	return &Installer{
		Config: config,
	}
}

// Parse installs parsed dependencies from a makefile.
func (i *Installer) Parse(r io.Reader) error {
	p := parser.New()

	nodes, err := p.Parse(r)
	if err != nil {
		return errors.Wrap(err, "parsing")
	}

	var paths []string
	for _, n := range nodes {
		switch n.(type) {
		case parser.Include:
			paths = append(paths, n.(parser.Include).Value)
		}
	}

	return i.Install(paths)
}

// Install the given dependencies.
func (i *Installer) Install(paths []string) error {
	done := make(chan struct{}, len(paths))
	errs := make(chan error, len(paths))

	for _, path := range paths {
		path := path
		go func() {
			err := i.installDependencyWithVariants(path)

			if err != nil {
				errs <- errors.Wrapf(err, "installing %q", path)
				return
			}

			done <- struct{}{}
		}()
	}

	for i := 0; i < len(paths); i++ {
		select {
		case err := <-errs:
			return err
		case <-done:
		}
	}

	return nil
}

// Install the dependency with lookup variants.
func (i *Installer) installDependencyWithVariants(path string) error {
	suffixes := []string{"", ".mk", "/index.mk"}

	for _, suffix := range suffixes {
		err := i.installDependency(path+suffix, path)

		// Unsupported (aka not github.com/* etc)
		if err == resolver.ErrNotSupported {
			return nil
		}

		// Remote not found, continue variants
		if err == resolver.ErrNotFound {
			continue
		}

		// Resolution error
		if err != nil {
			return err
		}

		// Found it
		return nil
	}

	return resolver.ErrNotFound
}

// Install dependency at the given path, stored using the original
// so that `include <path>` resolves properly. In the future it
// might be cool to support index.mk files and simply rewrite them.
func (i *Installer) installDependency(path, orig string) error {
	dir := filepath.Join(i.Destination, filepath.Dir(orig))
	dst := filepath.Join(i.Destination, orig)
	ctx := i.Log.WithField("path", path)

	// see if it exists
	if _, err := os.Stat(dst); err == nil {
		ctx.Debug("exists")
		return nil
	}

	// fetch it
	ctx.Info("install")
	r, err := i.Resolver.Get(path)
	if err != nil {
		return err
	}
	defer r.Close()

	// create the dir
	ctx.WithField("dir", dir).Debug("mkdir")
	if err := os.MkdirAll(dir, 0766); err != nil {
		return errors.Wrap(err, "mkdir")
	}

	// create the file
	ctx.WithField("dst", dst).Debug("copy")
	f, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "opening")
	}
	defer f.Close()

	// install its dependencies
	return i.Parse(io.TeeReader(r, f))
}
