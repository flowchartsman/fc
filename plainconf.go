package fc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
)

// PlainSource is a source for config files in an extremely simple format. Each
// line is tokenized as a single key/value pair. The first whitespace-delimited
// token in the line is interpreted as the flag name, and all remaining tokens
// are interpreted as the value. Any leading hyphens on the flag name are
// ignored.
type PlainSource struct {
	filename string
	m        map[string][]string
}

// WithConfigFile defines a new configuration source from the specified file
func WithConfigFile(filename string) *PlainSource {
	return &PlainSource{
		filename: filename,
	}
}

// Name returns a useful name for the plain config file source for usage
func (p *PlainSource) Name() string {
	return fmt.Sprintf("configuration file %q", p.filename)
}

// Loc simply returns the key where the value is expected to be found
func (p *PlainSource) Loc(key string) string {
	return fmt.Sprintf("%s, key %q", p.filename, key)
}

// Get returns the stringfied value stored at the specified key in the plain
// config file
func (p *PlainSource) Get(key string) ([]string, error) {
	if p.m == nil {
		if err := p.init(); err != nil {
			return nil, err
		}
	}
	values, ok := p.m[key]
	if !ok {
		return nil, ErrMissing
	}
	return values, nil
}

func (p *PlainSource) init() error {
	p.m = make(map[string][]string)

	cf, err := os.Open(p.filename)
	if err != nil {
		return err
	}
	defer cf.Close()

	s := bufio.NewScanner(cf)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue // skip empties
		}

		if line[0] == '#' {
			continue // skip comments
		}

		var (
			name  string
			value string
			index = strings.IndexRune(line, ' ')
		)
		if index < 0 {
			name, value = line, "true" // boolean option
		} else {
			name, value = line[:index], strings.TrimSpace(line[index:])
		}

		if i := strings.Index(value, " #"); i >= 0 {
			value = strings.TrimSpace(value[:i])
		}

		p.m[name] = strings.Split(value, ",")
	}
	return nil
}

// PlainFlagSource is a source that uses a flag value to define the config file to
// pull configuration from
type PlainFlagSource struct {
	*PlainSource
	flagName string
}

// WithConfigFileFlag defines a new configuration source from the configuration
// file provided by the specified flag
func WithConfigFileFlag(flag string) *PlainFlagSource {
	return &PlainFlagSource{
		PlainSource: &PlainSource{},
		flagName:    flag,
	}
}

// Name returns a useful name for the config file flag source for usage
func (pf *PlainFlagSource) Name() string {
	return fmt.Sprintf("configuration file defined by %q flag", pf.flagName)
}

// FlagNeeded returns the name of the flag that the PlainFlagSource will use to
// determine which file to pull configuration from
func (pf *PlainFlagSource) FlagNeeded() string {
	return pf.flagName
}

// WithFlagValue will will set the filename the PlainFlagSource will pull
// configuration from
func (pf *PlainFlagSource) WithFlagValue(value string) error {
	pf.PlainSource.filename = value
	if value == "" {
		return errors.New("PlainFlagSource given an empty string")
	}
	return nil
}
