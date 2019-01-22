package fc

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

// JSONSource is a source for config files in JSON format. Input should be
// an object. The object's keys are treated as flag names, and the object's
// values as flag values. If the value is an array, the flag will be set
// multiple times.
type JSONSource struct {
	filename string
	m        map[string]interface{}
}

// NewJSONSource defines a new configuration source from the specified JSON file
func NewJSONSource(filename string) *JSONSource {
	return &JSONSource{
		filename: filename,
	}
}

// Name returns a useful name for the JSON config source for usage
func (j *JSONSource) Name() string {
	return fmt.Sprintf("JSON configuration file %q", j.filename)
}

// Loc returns the object key where the value is expected to be found
func (j *JSONSource) Loc(key string) string {
	return fmt.Sprintf("%s, key %q", j.filename, key)
}

// Get returns the stringfied value stored at the specified key in the JSON file
func (j *JSONSource) Get(key string) ([]string, error) {
	if j.m == nil {
		if err := j.init(); err != nil {
			return nil, err
		}
	}
	values, err := stringifySlice(j.m[key])
	if err != nil {
		return nil, errors.Wrap(err, "error parsing JSON config")
	}
	return values, nil
}

func (j *JSONSource) init() error {
	m := make(map[string]interface{})

	jf, err := os.Open(j.filename)
	if err != nil {
		return err
	}
	defer jf.Close()

	d := json.NewDecoder(jf)
	// Must set UseNumber for stringifyValue to work
	d.UseNumber()
	err = d.Decode(&m)
	if err != nil {
		return errors.Wrap(err, "error parsing JSON config")
	}
	j.m = m
	return nil
}

func stringifySlice(val interface{}) ([]string, error) {
	if vals, ok := val.([]interface{}); ok {
		ss := make([]string, len(vals))
		for i := range vals {
			s, err := stringifyValue(vals[i])
			if err != nil {
				return nil, err
			}
			ss[i] = s
		}
		return ss, nil
	}
	s, err := stringifyValue(val)
	if err != nil {
		return nil, err
	}
	return []string{s}, nil
}

func stringifyValue(val interface{}) (string, error) {
	switch v := val.(type) {
	case string:
		return v, nil
	case json.Number:
		return v.String(), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", errors.Errorf("could not convert %q (type %T) to string", val, val)
	}
}
