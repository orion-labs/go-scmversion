# go-scmversion
Go (approximation) port of https://github.com/RiotGamesMinions/thor-scmversion

# Usage
```
Usage:
  go-scmversion [OPTIONS]

Application Options:
      --current  Print out the current version and end
      --auto     Bump the version based on what is found in the logs; default to #patch
      --major    Update Major version
      --minor    Update Minor version
      --patch    Update Patch version
      --pre=     Update prerelease
      --write    Actually write to git and output file
      --dir=     Directory from which to run the git commands
      --file=    File to write with the updated version number (default: ./VERSION)
      --debug    Enable debug logging of the version process

Help Options:
  -h, --help     Show this help message
```

# 3 Phases
A full process involves 3 phases of operation:
1. Calculate the current version based on the git tags
  1. If no acceptable version is found, it defaults to "0.0.0"
  2. You can limit the execution to just output this via the `--current` directive.
2. Increment the version forward the desired way
  1. The `--auto` directive looks for "#major"/"#minor"/"#patch" in the git logs since the Current tag, and will bump the appropriate level if found.
3. Store and output the updated version
  1. You must use the `--write` directive to make this phase happen
  2. This creates the tag locally, pushes tags upstream, and outputs to the `--file`.
