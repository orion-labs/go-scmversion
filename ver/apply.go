package ver

import (
	"fmt"
	"io"
	"os"

	"github.com/onbeep/go-scmversion/bump"
	"github.com/onbeep/go-scmversion/scm"
)

// Options are taken in from the command line
type Options struct {
	Current bool   `long:"current" description:"Print out the current version and end"`
	Auto    bool   `long:"auto" description:"Bump the version based on what is found in the logs; default to #patch"`
	Major   bool   `long:"major" description:"Update Major version"`
	Minor   bool   `long:"minor" description:"Update Minor version"`
	Patch   bool   `long:"patch" description:"Update Patch version"`
	Pre     string `long:"pre" description:"Update prerelease" default:""`
	Write   bool   `long:"write" description:"Actually write to git and output file"`
	Dir     string `long:"dir" description:"Directory from which to run the git commands"`
	File    string `long:"file" default:"./VERSION" description:"File to write with the updated version number"`
	Debug   bool   `long:"debug" description:"Enable debug logging of the version process"`
}

// NewProcessor builds the objects that perform the work
func NewProcessor(o Options) *Processor {
	g := scm.NewProvider(os.Stdout, o.Debug, o.Dir)
	return &Processor{Log: os.Stdout, Repo: g}
}

// Processor does the versioning work
type Processor struct {
	Log  io.Writer
	Repo scm.Provider
}

// Process actually does the work
func (p *Processor) Process(o *Options) error {
	current, err := p.Repo.Current()
	if err != nil {
		return err
	}
	fmt.Fprintf(p.Log, "Current Version: %s\n", current.String())
	if o.Current {
		return nil
	}

	var bumpMajor, bumpMinor, bumpPatch bool

	if o.Auto {
		bumpPatch = true
		bumpMajor, bumpMinor, err = p.Repo.Since(current)
	}

	updated := *current
	if bumpMajor || o.Major {
		updated, err = bump.Major(*current)
	} else if bumpMinor || o.Minor {
		updated, err = bump.Minor(*current)
	} else if bumpPatch || o.Patch {
		updated, err = bump.Patch(*current)
	} else if o.Pre != "" {
		updated, err = bump.Prerelease(*current, o.Pre)
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(p.Log, "Updated Version: %s\n", updated.String())
	if !o.Write {
		return nil
	}

	return p.Repo.Update(&updated)
}
