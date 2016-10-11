package scm

import (
	"os"
	"testing"

	"github.com/blang/semver"
)

func TestCandidates(t *testing.T) {
	dir := os.Getenv("THOR_DIR")
	if dir == "" {
		t.Skip("Missing env: THOR_DIR")
	}

	g := &gitter{Debug: true, Dir: dir}
	c, err := g.Current()
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}
	t.Logf("Current: %s\n", c.String())
}

func TestSince(t *testing.T) {
	dir := os.Getenv("THOR_DIR")
	if dir == "" {
		t.Skip("Missing env: THOR_DIR")
	}

	g := &gitter{Debug: true, Dir: dir}

	v, _ := semver.Make("0.0.9")
	major, minor, err := g.Since(&v)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}
	t.Logf("Major: %t; Minor: %t\n", major, minor)
}
