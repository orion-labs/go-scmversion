package scm

import (
	"os"
	"testing"
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
