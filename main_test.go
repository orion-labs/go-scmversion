package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/blang/semver"
	"github.com/onbeep/go-scmversion/cmd"
)

// OpenerFunc is a function that operats as an `Opener`
type OpenerFunc func(path string) (io.WriteCloser, error)

func (f OpenerFunc) Open(path string) (io.WriteCloser, error) {
	return f(path)
}

type closer struct {
	w io.Writer
}

func (c *closer) Write(b []byte) (n int, err error) {
	return c.w.Write(b)
}

func (c *closer) Close() error {
	return nil
}

func NopCloser(w io.Writer) io.WriteCloser {
	return &closer{w}
}

func TestCurrent(t *testing.T) {
	vStr := "9.8.7"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer

	o := &cmd.Options{Current: true}

	err := Process(o, &b, r, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	// t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) {
		t.Errorf("Missing output: %s\n", out)
	}
}

func TestPatch(t *testing.T) {
	vStr := "9.8.7"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer

	o := &cmd.Options{Patch: true, Debug: true}

	err := Process(o, &b, r, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) && !strings.Contains(out, "9.8.8") {
		t.Errorf("Missing output: %s\n", out)
	}
}

func TestMinor(t *testing.T) {
	vStr := "9.8.7"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer

	o := &cmd.Options{Minor: true, Patch: true, Debug: true}

	err := Process(o, &b, r, nil)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) && !strings.Contains(out, "9.9.0") {
		t.Errorf("Missing output: %s\n", out)
	}
}

func TestMajor(t *testing.T) {
	vStr := "9.8.7"
	end := "10.0.0"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer

	file := "./bubbagump.txt"
	o := &cmd.Options{Major: true, Minor: true, Patch: true, Debug: true, Write: true, File: file}

	var d bytes.Buffer
	dest := NopCloser(&d)
	op := func(path string) (io.WriteCloser, error) { return dest, nil }
	err := Process(o, &b, r, OpenerFunc(op))
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) && !strings.Contains(out, end) {
		t.Errorf("Missing output: %s\n", out)
	}

	written := d.String()
	t.Logf("Written: %s\n", written)
	if !strings.Contains(written, end) {
		t.Errorf("Missing write: %s\n", written)
	}
}

type MockRepo struct {
	Ver        *semver.Version
	HasMajor   bool
	HasMinor   bool
	UpdateFail bool
}

func (m *MockRepo) Current() (*semver.Version, error) {
	if m.Ver == nil {
		return nil, fmt.Errorf("forced error")
	}
	return m.Ver, nil
}

func (m *MockRepo) Since(v *semver.Version) (bool, bool, error) {
	return m.HasMajor, m.HasMinor, nil
}

func (m *MockRepo) Update(v *semver.Version) error {
	if m.UpdateFail {
		return fmt.Errorf("forced error")
	}
	return nil
}
