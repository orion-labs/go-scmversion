package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/orion-labs/go-scmversion/bump"
	"github.com/orion-labs/go-scmversion/cmd"
	"github.com/orion-labs/go-scmversion/scm"
)

// Opener is a factory to open the local destination of the version number
type Opener interface {
	Open(path string) (io.WriteCloser, error)
}

type fileOpen struct{}

func (o *fileOpen) Open(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

// Process actually does the work
func Process(o *cmd.Options, log io.Writer, repo scm.Provider, op Opener) error {
	current, err := repo.Current()
	if err != nil {
		return err
	}
	fmt.Fprintf(log, "Current Version: %s\n", current.String())
	if o.Current {
		return nil
	}

	var bumpMajor, bumpMinor, bumpPatch bool

	if o.Auto {
		bumpPatch = true
		bumpMajor, bumpMinor, err = repo.Since(current)
	}

	updated := *current
	if bumpMajor || o.Major {
		if o.Debug {
			fmt.Fprintf(log, "Bump Major\n")
		}
		updated, err = bump.Major(*current)
	} else if bumpMinor || o.Minor {
		if o.Debug {
			fmt.Fprintf(log, "Bump Minor\n")
		}
		updated, err = bump.Minor(*current)
	} else if bumpPatch || o.Patch {
		if o.Debug {
			fmt.Fprintf(log, "Bump Patch\n")
		}
		updated, err = bump.Patch(*current)
	} else if o.Pre != "" {
		if o.Debug {
			fmt.Fprintf(log, "Bump Prerelease: %s\n", o.Pre)
		}
		updated, err = bump.Prerelease(*current, o.Pre)
	}
	if err != nil {
		return err
	}

	r := &updated
	if r.Equals(*current) {
		fmt.Fprintf(log, "No update found: %s\n", updated.String())
		return fmt.Errorf("no update found")
	}

	fmt.Fprintf(log, "Updated Version: %s\n", updated.String())
	if !o.Write {
		fmt.Fprintf(log, "No write requested\n")
		return nil
	}

	file, err := op.Open(o.File)
	if err != nil {
		return err
	}
	defer file.Close()
	fmt.Fprintf(file, "%s", updated.String())

	err = repo.Update(&updated)
	if err != nil {
		fmt.Fprintf(log, "Update error\n")
	}
	return err
}

func main() {
	var options cmd.Options
	parser := flags.NewParser(&options, flags.Default)
	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	var log io.Writer = os.Stdout
	scm := scm.NewProvider(log, &options)

	err := Process(&options, log, scm, &fileOpen{})
	exit := 0
	if err != nil {
		fmt.Fprintf(log, "Exit w/error: %v\n", err)
		exit = 1
	}
	os.Exit(exit)
}
