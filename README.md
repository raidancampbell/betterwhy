# betterwhy
A better version of 'go mod why' which explains build requirements not actually required for execution

## Usage

```
betterwhy -h                                                                                                                                                                                                [12:26:30]
betterwhy traverses go mod's dependency tree to diagnose why a given dependency is required.  
It is different from 'go mod why' in that it will look at dependencies that are required to be pulled in, but not required for the build.
This is useful for scenarios like:
 - a direct dependency's test dependencies (no longer a problem in 1.16+)
 - explain why an exact version of a dependency was pulled in, allowing you to explain every entry in go.sum
It does not handle indirect dependencies (e.g. any dependency that does not declare its own dependencies, but still requires them).
Additionally, betterwhy finds A path to the desired dependency.
If two dependencies ultimately pull in the desired version, the returned path isn't guaranteed to be the shortest one.'

Usage:
  betterwhy <dependency with version to explain>

Examples:
betterwhy github.com/stretchr/testify@v1.3.0

Flags:
  -h, --help   help for betterwhy
```

## Example
```shell
go install
``` 
and 
```shell
betterwhy github.com/stretchr/testify@v1.3.0
```
Sample output for this repository:
```
betterwhy github.com/stretchr/testify@v1.3.0                                                                                                                                                                [12:25:11]
# github.com/stretchr/testify@v1.3.0
github.com/raidancampbell/betterwhy
github.com/spf13/cobra@v1.1.3
github.com/spf13/viper@v1.7.0
github.com/bketelsen/crypt@v0.0.3-0.20200106085610-5cbc8cc4026c
github.com/hashicorp/consul/api@v1.1.0
github.com/hashicorp/serf@v0.8.2
github.com/stretchr/testify@v1.3.0
```