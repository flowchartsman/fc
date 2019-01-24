# fc [![Latest Release](https://img.shields.io/github/release/flowchartsman/fc.svg?style=flat-square)](https://github.com/flowchartsman/fc/releases/latest) [![GoDoc](https://godoc.org/github.com/flowchartsman/fc?status.svg)](https://godoc.org/github.com/flowchartsman/fc) [![Travis CI](https://travis-ci.org/flowchartsman/fc.svg?branch=master)](https://travis-ci.org/flowchartsman/fc)

fc ("FlagConf") is a flexible configuration library using the stdlib's
[flag.FlagSet](https://golang.org/pkg/flag#FlagSet) as its source of
definitions
([rationale](https://peter.bourgon.org/go-for-industrial-programming/#program-configuration)).
It's based on [Peter Bourgon's ff
library](https://github.com/peterbourgon/ff/blob/master/README.md), but allows
you to plug in any source of configuration you like, in any order.

## Usage

Define a flag.FlagSet in your func main.

```go
func main() {
	fs := flag.NewFlagSet("my-program", flag.ExitOnError)
	var (
		listenAddr = fs.String("listen-addr", "localhost:8080", "listen address")
		refresh    = fs.Duration("refresh", 15*time.Second, "refresh interval")
		debug      = fs.Bool("debug", false, "log debug information")
	)
```

Then, call fc.Parse instead of fs.Parse.

```go
   fc.Parse(fs,
        fc.WithEnv("MY_PROGRAM"),
        fc.WithConfigFile("myprogram.conf"),
   )
```

This example will parse flags from the commandline args, just like regular
package flag, with the highest priority. Then it will look in the environment
for variables with a `MY_PROGRAM` prefix. After that, if it still hasn't found
all of the values it needs, it will attmept to load the values from
`myprogram.conf` which expects a file in this format:

```
listen-addr localhost:8080
refresh 30s
debug true
```
## JSON files
`WithJSONFile(filename string)` is also provided to decode JSON files, which
should be a single JSON object, with keys corresponding to the flag names:

```json
{
    "listen-addr":"localhost:8080",
    "refresh":"30s",
    "debug":true
}
```

## Adding configuration sources
In order to keep non-stdlib dependencies to a minimum, other configurations
sources should exist in other repositories, and should provide types which
correspond to the `fc.Source` interface:

```go
type Source interface {
	// Get attempts to retrieve any values for the given key from the store
	Get(key string) ([]string, error)
	// Name prints a meaningful name for the configuration store for the usage
	// message
	Name() string
	// Loc prints a meaningful store location name for the key for the usage
	// message
	Loc(key string) string
}
```

`Name()` and `Loc()` are intended to be used to generate informative error
messages and to provide an extra stanza at the end of the usage message,
detailing the configuration sources that will be pulled from.

## Other sources

Some additional configurations sources (more will be added as they are created):

- [consul](https://github.com/flowchartsman/fc-consul)

## Guidelines
It goes without saying that the more configuration sources you use, the more
difficult it is to determine where values are coming from, and the greater the
possibility your program will fail at runtime, so use your best judgment.

You should almost always be parsing from the local ENV before any other sources
of configuration, in keeping with
[12-factor app guidelines](https://12factor.net/config).

As well, sources should be lazily-loaded so that your program doesn't pull in
any extra data that it doesn't need during startup. Look to `plainconf.go` and
`jsonconf.go` for guidelines.
