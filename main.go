package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/onbeep/go-scmversion/ver"
)

var options ver.Options
var parser = flags.NewParser(&options, flags.Default)

func main() {
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}
	p := ver.NewProcessor(options)
	err := p.Process(&options)
	exit := 0
	if err != nil {
		exit = 1
	}
	os.Exit(exit)
}
