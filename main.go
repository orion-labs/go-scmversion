package main

import (
	"io"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/onbeep/go-scmversion/cmd"
	"github.com/onbeep/go-scmversion/scm"
	"github.com/onbeep/go-scmversion/ver"
)

var options cmd.Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
	var log io.Writer = os.Stdout
	scm := scm.NewProvider(log, &options)
	p := ver.NewProcessor(log, scm)

	err := p.Process(&options)
	exit := 0
	if err != nil {
		exit = 1
	}
	os.Exit(exit)
}
