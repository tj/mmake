package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	log "github.com/apex/log"
	"github.com/apex/log/handlers/cli"

	"github.com/tj/mmake/export"
	"github.com/tj/mmake/help"
	"github.com/tj/mmake/installer"
	"github.com/tj/mmake/resolver"
)

var makefileVariants = [3]string{"GNUmakefile", "Makefile", "makefile"}

func init() {
	var level = os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "warning"
	}
	log.SetHandler(cli.Default)
	log.SetLevel(log.MustParseLevel(level))
}

func main() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	f, err := findMakefile()
	if err != nil {
		log.WithError(err).Fatal("Unable to find makefile")
	}

	b, err := ioutil.ReadFile(f)
	if err != nil {
		log.WithError(err).Fatal("reading makefile")
	}

	// ensure deps are installed
	i := installer.New(installer.Config{
		Resolver:    resolver.NewUniversalResolver(),
		Destination: resolver.GetIncludePath(os.Args[1:]),
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
	case "export":
		s, _ := export.Export(f, bytes.NewReader(b))
		fmt.Println(s)
	default:
		passThrough(os.Args[1:])
	}
}

// doHelp outputs target comments.
func doHelp(r io.Reader, args []string) {
	var err error

	if len(args) == 0 {
		err = help.OutputAllShort(r, os.Stdout, []string{})
	} else if args[0] == "-v" || args[0] == "--verbose" {
		err = help.OutputAllLong(r, os.Stdout, args[1:])
	} else {
		err = help.OutputAllShort(r, os.Stdout, args)
	}

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

// Follows the manual from GNUMake which will attempt the following
// three filenames until successfully reading the first one:
// GNUmakefile, makefile and Makefile
// https://www.gnu.org/software/make/manual/make.html#Makefile-Names
// GNUmakefile, makefile and Makefile
func findMakefile() (string, error) {
	for _, m := range makefileVariants {
		if fileExists(m) {
			return m, nil
		}
	}
	return "", errors.New("No makefile found among in: " + fmt.Sprint(makefileVariants))
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
