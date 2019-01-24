package fc

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// ErrMissing is returned when the given key is not in a store
var ErrMissing = errors.New("key not present in the store")

// Source represents a source of configuration information
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

// ParseArgs parses the provided arguments with the given FlagSet and sources,
// starting with the commandline flags and progressing through all given
// sources in decreasing priority order until a value is found
func ParseArgs(args []string, fs *flag.FlagSet, sources ...Source) error {
	fs.Usage = fcUsage(fs, sources)
	err := fs.Parse(args)
	if err != nil {
		return err
	}
	found := map[string]bool{}
	fs.Visit(func(f *flag.Flag) {
		found[f.Name] = true
	})

	fs.VisitAll(func(f *flag.Flag) {
		// Bail if we've encountered an error
		if err != nil {
			return
		}
		for _, source := range sources {
			if found[f.Name] {
				return
			}
			values, serr := source.Get(f.Name)
			switch serr {
			case ErrMissing:
				continue
			case nil:
			default:
				err = serr
				return
			}
			for _, individual := range values {
				if verr := fs.Set(f.Name, strings.TrimSpace(individual)); verr != nil {
					err = errors.Wrapf(verr, "error setting flag %q from %s", f.Name, source.Loc(f.Name))
					return
				}
				found[f.Name] = true
			}
		}
	})

	return err
}

// Parse is a convenient alias for ParseArgs targeting os.Args[1:]
func Parse(fs *flag.FlagSet, sources ...Source) error {
	return ParseArgs(os.Args[1:], fs, sources...)
}

func fcUsage(fs *flag.FlagSet, sources []Source) func() {
	return func() {
		if fs.Name() == "" {
			fmt.Fprintf(fs.Output(), "Usage:\n")
		} else {
			fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
		}
		fs.PrintDefaults()
		if len(sources) > 0 {
			fmt.Fprintln(fs.Output(), "\nAdditional configuration sources:")
			for _, s := range sources {
				fmt.Fprintf(fs.Output(), "\t- %s\n", s.Name())
			}
		}
	}
}
