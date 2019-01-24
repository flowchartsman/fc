package fc

import (
	"fmt"
	"os"
	"strings"
)

// EnvSource represents a config source from the environment with an optional
// prefix. By convention, the prefix and flag name will be converted to
// UPPERCASE, and all dashes will be converted to underscores.
type EnvSource struct {
	prefix string
}

// WithEnv returns a new source that pulls from os.ENV with the specified
// prefix, or from the entire environment if no prefix is provided
func WithEnv(prefix string) *EnvSource {
	return &EnvSource{prefix: strings.ToUpper(prefix) + "_"}
}

// Get attempts to retrieve a flag from os.ENV
func (e *EnvSource) Get(key string) ([]string, error) {
	key = e.Loc(key)
	value, found := os.LookupEnv(key)
	if !found {
		return nil, ErrMissing
	}
	return strings.Split(value, ","), nil
}

// Loc returns the computed environment name for the flag
func (e *EnvSource) Loc(key string) string {
	key = strings.ToUpper(key)
	if e.prefix != "" {
		key = e.prefix + key
	}
	return envVarReplacer.Replace(key)
}

// Name returns a useful name for the EnvSource for usage
func (e *EnvSource) Name() string {
	if e.prefix == "" {
		return "environment variables"
	}
	return fmt.Sprintf("environment variables with the prefix %q", e.prefix)
}

var envVarReplacer = strings.NewReplacer(
	"-", "_",
	".", "_",
	"/", "_",
)
