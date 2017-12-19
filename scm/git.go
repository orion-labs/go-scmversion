package scm

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/blang/semver"
	"github.com/onbeep/go-scmversion/cmd"
)

// Provider is the abstraction over interacting with an actual SCM
type Provider interface {
	Current() (*semver.Version, error)
	Since(v *semver.Version) (hasMajor bool, hasMinor bool, err error)
	Update(v *semver.Version) error
}

// NewProvider creates the source-code management tool
func NewProvider(log io.Writer, o *cmd.Options) Provider {
	return &gitter{Log: log, Debug: o.Debug, Dir: o.Dir}
}

type gitter struct {
	Log   io.Writer
	Debug bool
	Dir   string
}

func (g *gitter) cmd(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	if g.Dir != "" {
		cmd.Dir = g.Dir
	}
	return cmd
}

// Current retrieves the latest version for the current git working directory
func (g *gitter) Current() (*semver.Version, error) {
	// First, ensure we have all the info locally
	fetch := g.cmd("fetch", "--all")
	content, err := fetch.Output()
	if err != nil {
		return nil, err
	}
	if g.Debug {
		fmt.Printf("Fetch output: %s\n", content)
	}

	// Grab all the tags for the repo
	tag := g.cmd("tag")
	var tagOut bytes.Buffer
	tag.Stdout = &tagOut
	err = tag.Run()
	if err != nil {
		return nil, err
	}

	// Cycle thru the tags, and evaluate if each might be the latest
	r, _ := semver.Make("0.0.0")
	parts := strings.Split(tagOut.String(), "\n")
	for ix := len(parts) - 1; ix >= 0; ix-- {
		tag := parts[ix]
		if tag == "" {
			continue
		}
		branch := g.cmd("branch", "--contains", tag)
		var branchOut bytes.Buffer
		branch.Stdout = &branchOut
		berr := branch.Run()
		if berr != nil {
			fmt.Printf("Err: %s - %v\n", tag, berr)
		}
		if g.Debug {
			fmt.Fprintf(g.Log, "Tag: %s\n", tag)
		}

		c, cerr := semver.Make(tag)
		if cerr != nil {
			if g.Debug {
				fmt.Fprintf(g.Log, "Format err: %s\n", tag)
			}
			continue
		}

		if !c.GT(r) {
			continue
		}

		// Tring to simulate this: https://github.com/RiotGamesMinions/thor-scmversion/blob/4cf603205075aa7ed8cb3d772fa6b006fcce09fc/lib/thor-scmversion/git_version.rb
		// def contained_in_current_branch?(tag)
		//   ShellUtils.sh("git branch --contains #{tag}") =~ /\*/
		// end
		found := false
		containers := strings.Split(branchOut.String(), "\n")
		for _, contain := range containers {
			if strings.HasPrefix(contain, "*") {
				found = true
				break
			}
		}
		if !found {
			if g.Debug {
				fmt.Fprintf(g.Log, "Not in branch: %s\n", tag)
			}
			continue
		}
		r = c
	}

	return &r, nil
}

// Singe the given tag, fetch the short version of the logs
func (g *gitter) Since(v *semver.Version) (hasMajor bool, hasMinor bool, err error) {
	since := fmt.Sprintf("%s..", v.String())
	search := g.cmd("log", "--abbrev-commit", "--format=oneline", since)

	logs, err := search.Output()
	if err != nil {
		return false, false, err
	}
	if g.Debug {
		fmt.Printf("Fetch output: %s\n", string(logs))
	}
	major := bytes.Contains(logs, []byte("#major"))
	minor := bytes.Contains(logs, []byte("#minor"))
	return major, minor, nil
}

// Update the repository (and upstream) with the given version as a tag
func (g *gitter) Update(v *semver.Version) error {
	tag := v.String()
	cmt := fmt.Sprintf("Version %s", tag)

	apply := g.cmd("tag", "-a", "-m", cmt, tag)
	err := apply.Run()
	if err != nil {
		return err
	}

	push := g.cmd("push", "--tags")
	return push.Run()
}
