package ver

import (
	"fmt"
	"io"

	"github.com/onbeep/go-scmversion/bump"
	"github.com/onbeep/go-scmversion/cmd"
	"github.com/onbeep/go-scmversion/scm"
)

// NewProcessor builds the objects that perform the work
func NewProcessor(log io.Writer, p scm.Provider) *Processor {
	return &Processor{Log: log, Repo: p}
}

// Processor does the versioning work
type Processor struct {
	Log  io.Writer
	Repo scm.Provider
}

// Process actually does the work
// TODO: this can maybe just be turned into a function
func (p *Processor) Process(o *cmd.Options) error {
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
		if o.Debug {
			fmt.Fprintf(p.Log, "Bump Major\n")
		}
		updated, err = bump.Major(*current)
	} else if bumpMinor || o.Minor {
		if o.Debug {
			fmt.Fprintf(p.Log, "Bump Minor\n")
		}
		updated, err = bump.Minor(*current)
	} else if bumpPatch || o.Patch {
		if o.Debug {
			fmt.Fprintf(p.Log, "Bump Patch\n")
		}
		updated, err = bump.Patch(*current)
	} else if o.Pre != "" {
		if o.Debug {
			fmt.Fprintf(p.Log, "Bump Prerelease: %s\n", o.Pre)
		}
		updated, err = bump.Prerelease(*current, o.Pre)
	}
	if err != nil {
		return err
	}

	r := &updated
	if r.Equals(*current) {
		fmt.Fprintf(p.Log, "No update found: %s\n", updated.String())
		return fmt.Errorf("no update found")
	}

	fmt.Fprintf(p.Log, "Updated Version: %s\n", updated.String())
	if !o.Write {
		return nil
	}

	return p.Repo.Update(&updated)
}
