package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/segmentio/go-env"

	"github.com/tj/mmake/help"
	"github.com/tj/mmake/installer"
	"github.com/tj/mmake/resolver"
)

func init() {
	log.SetHandler(cli.Default)

	level, err := log.ParseLevel(env.GetDefault("LOG_LEVEL", "warning"))
	if err != nil {
		panic("invalid log level")
	}

	log.SetLevel(level)
}

func main() {
	// read Makefile
	b, err := ioutil.ReadFile("Makefile")
	if err != nil {
		log.WithError(err).Fatal("reading makefile")
	}

	// ensure deps are installed
	i := installer.New(installer.Config{
		Resolver:    resolver.NewGithubResolver(),
		Destination: "/usr/local/include",
		Log:         log.Log,
	})

	if err := i.Parse(bytes.NewReader(b)); err != nil {
		log.WithError(err).Fatal("installing")
	}

	// output target help
	if len(os.Args) > 2 && os.Args[1] == "help" {
		err := help.OutputTargetLong(bytes.NewReader(b), os.Stdout, os.Args[2])
		if err != nil {
			log.WithError(err).Fatal("outputting help")
		}
		return
	}

	// output all help
	if len(os.Args) > 1 && os.Args[1] == "help" {
		err := help.OutputAllShort(bytes.NewReader(b), os.Stdout)
		if err != nil {
			log.WithError(err).Fatal("outputting help")
		}
		return
	}

	// make pass-through
	cmd := exec.Command("make", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()

	if e, ok := err.(*exec.ExitError); ok {
		if status, ok := e.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		} else {
			os.Exit(1)
		}
	}

	if err != nil {
		log.WithError(err).Fatal("executing")
	}
}
