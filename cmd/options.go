package cmd

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
