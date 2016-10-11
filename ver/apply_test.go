package ver

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/blang/semver"
)

func TestCurrent(t *testing.T) {
	vStr := "9.8.7"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer
	p := &processor{Log: &b, Repo: r}

	o := &Options{Current: true}

	err := p.Process(o)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	// t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) {
		t.Errorf("Missing output: %s\n", out)
	}
}

func TestMajor(t *testing.T) {
	vStr := "9.8.7"
	v, _ := semver.Make(vStr)
	r := &MockRepo{Ver: &v}
	var b bytes.Buffer
	p := &processor{Log: &b, Repo: r}

	o := &Options{DryRun: true, Major: true}

	err := p.Process(o)
	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	out := b.String()
	t.Logf("Output: %s\n", out)
	if !strings.Contains(out, vStr) && !strings.Contains(out, "10.0.0") {
		t.Errorf("Missing output: %s\n", out)
	}
}

type MockRepo struct {
	Ver        *semver.Version
	Logs       bytes.Buffer
	UpdateFail bool
}

func (m *MockRepo) Current() (*semver.Version, error) {
	if m.Ver == nil {
		return nil, fmt.Errorf("forced error")
	}
	return m.Ver, nil
}

func (m *MockRepo) Since(v *semver.Version) ([]byte, error) {
	return m.Logs.Bytes(), nil
}

func (m *MockRepo) Update(v *semver.Version) error {
	if m.UpdateFail {
		return fmt.Errorf("forced error")
	}
	return nil
}
