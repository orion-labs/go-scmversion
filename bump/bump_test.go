package bump

import (
	"testing"

	"github.com/blang/semver"
)

func TestMajor(t *testing.T) {
	var majorData = []struct {
		in  string
		out string
	}{
		{"0.0.0", "1.0.0"},
		{"2.3.4", "3.0.0"},
		{"4.5.6+build", "5.0.0"},
		{"6.7.8-alpha", "7.0.0"},
		{"8.9.10-beta+arbitrary", "9.0.0"},
	}

	for _, tc := range majorData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Major(s)
			if err != nil {
				t.Fatalf("Problem bumping: %s -> %s\n", tc.in, err)
			}
			r := n.String()
			if r != tc.out {
				t.Errorf("Unexpected result: %s -> %s != %s\n", tc.in, r, tc.out)
			}
		})
	}
}

func TestMinor(t *testing.T) {
	var minorData = []struct {
		in  string
		out string
	}{
		{"0.0.0", "0.1.0"},
		{"2.3.4", "2.4.0"},
		{"4.5.6+build", "4.6.0"},
		{"6.7.8-alpha", "6.8.0"},
		{"8.9.10-beta+arbitrary", "8.10.0"},
	}

	for _, tc := range minorData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Minor(s)
			if err != nil {
				t.Fatalf("Problem bumping: %s -> %s\n", tc.in, err)
			}
			r := n.String()
			if r != tc.out {
				t.Errorf("Unexpected result: %s -> %s != %s\n", tc.in, r, tc.out)
			}
		})
	}
}

func TestPatch(t *testing.T) {
	var patchData = []struct {
		in  string
		out string
	}{
		{"0.0.0", "0.0.1"},
		{"2.3.4", "2.3.5"},
		{"4.5.6+build", "4.5.7"},
		{"6.7.8-alpha", "6.7.9"},
		{"8.9.10-beta+arbitrary", "8.9.11"},
	}

	for _, tc := range patchData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Patch(s)
			if err != nil {
				t.Fatalf("Problem bumping: %s -> %s\n", tc.in, err)
			}
			r := n.String()
			if r != tc.out {
				t.Errorf("Unexpected result: %s -> %s != %s\n", tc.in, r, tc.out)
			}
		})
	}
}

func TestBuildSuccess(t *testing.T) {
	buildVal := "EXAMPLE.BUILD"
	var buildData = []struct {
		in  string
		out string
	}{
		{"0.0.0", "0.0.0+EXAMPLE.BUILD"},
		{"2.3.4", "2.3.4+EXAMPLE.BUILD"},
		{"4.5.6+build", "4.5.6+EXAMPLE.BUILD"},
		{"6.7.8-alpha", "6.7.8-alpha+EXAMPLE.BUILD"},
		{"8.9.10-beta+arbitrary.addition", "8.9.10-beta+EXAMPLE.BUILD"},
	}

	for _, tc := range buildData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Build(s, buildVal)
			if err != nil {
				t.Fatalf("Problem bumping: %s -> %s\n", tc.in, err)
			}
			r := n.String()
			if r != tc.out {
				t.Errorf("Unexpected result: %s -> %s != %s\n", tc.in, r, tc.out)
			}
		})
	}
}

func TestBuildFailure(t *testing.T) {
	buildErrs := []string{
		"no_unders",
		"multi..dot",
		"other?chars",
	}
	start, _ := semver.Make("1.2.3")
	for _, build := range buildErrs {
		t.Run(build, func(t *testing.T) {
			n, err := Build(start, build)
			if err == nil {
				t.Fatalf("Expected error: %s -> %s\n", build, n.String())
			}
		})
	}
}

func TestPrereleaseSuccess(t *testing.T) {
	var preData = []struct {
		in  string
		pre string
		out string
	}{
		{"0.0.0", "apple", "0.0.1-apple"},
		{"1.1.1-baker", "apple", "1.1.1-apple"},
		{"1.2.3-granny.smith", "apple", "1.2.3-apple"},
		{"1.3.5-caramel.apple.stick", "apple", "1.3.5-apple"},
		{"2.3.4-pear", "pear", "2.3.4-pear.1"},
		{"3.4.5-pear.9", "pear", "3.4.5-pear.10"},
		{"3.6.9-pear.99.bottles", "pear", "3.6.9-pear.100"},
		{"4.5.6", "99", "4.5.7-99"},
		{"4.5.7-99", "100", "4.5.7-100"},
		{"4.5.7-100", "100", "4.5.7-100.1"},
	}

	for _, tc := range preData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Prerelease(s, tc.pre)
			if err != nil {
				t.Fatalf("Problem bumping: %s + %s -> %s\n", tc.in, tc.pre, err)
			}
			r := n.String()
			if r != tc.out {
				t.Errorf("Unexpected result: %s + %s -> %s != %s\n", tc.in, tc.pre, r, tc.out)
			}
		})
	}
}

func TestPrereleaseError(t *testing.T) {
	var preData = []struct {
		in  string
		pre string
	}{
		{"0.0.0", "xray.zulu"},
		{"1.22.333", ""},
	}

	for _, tc := range preData {
		t.Run(tc.in, func(t *testing.T) {
			s, err := semver.Make(tc.in)
			if err != nil {
				t.Fatalf("Invalid input: %s -> %s\n", tc.in, err)
			}
			n, err := Prerelease(s, tc.pre)
			if err == nil {
				t.Fatalf("Expected error: %s + %s -> %s\n", tc.in, tc.pre, n.String())
			}
		})
	}
}
