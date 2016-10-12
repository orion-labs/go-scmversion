package bump

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blang/semver"
)

// Major is for API incompatible changes
func Major(v semver.Version) (semver.Version, error) {
	r := semver.Version{
		Major: (v.Major + 1),
	}
	return r, nil
}

// Minor is for when you add functionality in a backwards-compatible manner
func Minor(v semver.Version) (semver.Version, error) {
	r := semver.Version{
		Major: v.Major,
		Minor: (v.Minor + 1),
	}
	return r, nil
}

// Patch is for when you make backwards-compatible bug fixes
func Patch(v semver.Version) (semver.Version, error) {
	r := semver.Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: (v.Patch + 1),
	}
	return r, nil
}

// Build is for metadata without precedence
func Build(v semver.Version, build string) (semver.Version, error) {
	base := strings.Split(v.String(), "+")[0]
	return semver.Parse(strings.Join([]string{base, build}, "+"))
}

// Prerelease is for versions have a lower precedence than the associated normal version
func Prerelease(v semver.Version, pre string) (semver.Version, error) {
	if len(strings.Split(pre, ".")) > 1 {
		return semver.Version{}, fmt.Errorf("multi-field prerelease not accepted")
	}

	// Build the Prerelease object
	pv, err := semver.NewPRVersion(pre)
	if err != nil {
		return semver.Version{}, err
	}
	// if pv.IsNum {
	// 	return semver.Version{}, fmt.Errorf("bump non-numeric prerelease only")
	// }

	if len(v.Pre) == 0 {
		p, perr := Patch(v)
		if perr != nil {
			return p, perr
		}
		p.Pre = []semver.PRVersion{pv}
		return p, nil
	}

	r := semver.Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch,
	}

	if v.Pre[0].VersionStr != pv.VersionStr || v.Pre[0].VersionNum != pv.VersionNum {
		r.Pre = []semver.PRVersion{pv}
		return r, nil
	}
	var count uint64
	if len(v.Pre) > 1 && v.Pre[1].IsNum {
		count = v.Pre[1].VersionNum
	}
	count = count + 1
	sv, err := semver.NewPRVersion(strconv.FormatUint(count, 10))
	if err != nil {
		return semver.Version{}, fmt.Errorf("bump problem with prerelease")
	}
	r.Pre = []semver.PRVersion{pv, sv}
	return r, nil
}
