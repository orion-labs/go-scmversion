package scm

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/blang/semver"
)

type gitter struct {
	Debug bool
	Dir   string
}

// Current retrieves the latest version for the current git working directory
func (g *gitter) Current() (*semver.Version, error) {
	// First, ensure we have all the info locally
	fetch := exec.Command("git", "fetch", "--all")
	if g.Dir != "" {
		fetch.Dir = g.Dir
	}
	content, err := fetch.Output()
	if err != nil {
		return nil, err
	}
	if g.Debug {
		fmt.Printf("Fetch output: %s\n", content)
	}

	// Grab all the tags for the repo
	tag := exec.Command("git", "tag")
	if g.Dir != "" {
		tag.Dir = g.Dir
	}
	var tagOut bytes.Buffer
	tag.Stdout = &tagOut
	err = tag.Run()
	if err != nil {
		return nil, err
	}

	// Cycle thru the tags, and evaluate if each might be the latest
	r, _ := semver.Make("0.0.0")
	parts := strings.Split(tagOut.String(), "\n")
	// candidates := make([]string, 0, len(parts))
	for _, tag := range parts {
		if tag == "" {
			continue
		}
		branch := exec.Command("git", "branch", "--contains", tag)
		if g.Dir != "" {
			branch.Dir = g.Dir
		}
		berr := branch.Run()
		if berr != nil {
			fmt.Printf("Err: %s - %v\n", tag, berr)
		}
		fmt.Printf("Tag: %s\n", tag)

		c, cerr := semver.Make(tag)
		if cerr != nil {
			fmt.Printf("Format err: %s\n", tag)
			continue
		}

		if c.GT(r) {
			r = c
		}
	}

	return &r, nil
}

// Update the repository (and upstream) with the given version as a tag
func (g *gitter) Update(v *semver.Version) error {
	tag := v.String()
	cmt := fmt.Sprintf("Version %s", tag)

	apply := exec.Command("git", "-a", "-m", cmt, tag)
	if g.Dir != "" {
		apply.Dir = g.Dir
	}
	err := apply.Run()
	if err != nil {
		return err
	}

	push := exec.Command("git", "push", "--tags")
	if g.Dir != "" {
		push.Dir = g.Dir
	}
	return push.Run()
}
