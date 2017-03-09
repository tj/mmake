package main

import (
	"bytes"
	"io"
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
	log.SetLevel(log.MustParseLevel(env.GetDefault("LOG_LEVEL", "warning")))
}

func main() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

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
		ForceUpdate: cmd == "update",
	})

	if err := i.Parse(bytes.NewReader(b)); err != nil {
		log.WithError(err).Fatal("installing")
	}

	switch cmd {
	case "update":
		return
	case "help":
		doHelp(bytes.NewReader(b), os.Args[2:])
	default:
		passThrough(os.Args[1:])
	}
}

// doHelp outputs target comments.
func doHelp(r io.Reader, args []string) {
	if len(args) > 0 {
		var err error

		if args[0] == "-v" {
			err = help.OutputAllLong(r, os.Stdout)
		} else {
			err = help.OutputTargetLong(r, os.Stdout, args[0])
		}

		if err != nil {
			log.WithError(err).Fatal("outputting help")
		}
		return
	}

	// output all help
	err := help.OutputAllShort(r, os.Stdout)
	if err != nil {
		log.WithError(err).Fatal("outputting help")
	}
}

// passThrough executes make.
func passThrough(args []string) {
	cmd := exec.Command("make", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

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
