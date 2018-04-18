package otto

import (
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// Params holds url params and provides
// simple ways to parse the value to different types
type Params map[string]string

// String returns the value as string
func (p Params) String(key string) string {
	return p[key]
}

// Int returns the value as int
func (p Params) Int(key string) (int, error) {
	vi, err := strconv.Atoi(p[key])
	return vi, errors.Wrapf(err, "failed to parse '%v' to int", p[key])
}

// Bool returns the value as bool
func (p Params) Bool(key string) (bool, error) {
	vb, err := strconv.ParseBool(p[key])
	return vb, errors.Wrapf(err, "failed to parse '%v' to bool", p[key])
}

// ValueParams holds url.Values and provides
// simple ways to parse the value to different types
type ValueParams struct {
	vals url.Values
}

// String gets one value associated with key and returns it as string
func (p ValueParams) String(key string) string {
	return p.vals.Get(key)
}

// Strings returns all value associated with key as string slice
func (p ValueParams) Strings(key string) []string {
	return p.vals[key]
}

// Int gets one value associated with key and returns it as int
func (p ValueParams) Int(key string) (int, error) {
	v := p.vals.Get(key)
	vi, err := strconv.Atoi(v)
	return vi, errors.Wrapf(err, "failed to parse '%v' to int", v)
}

// Ints returns all value associated with key as int slice
func (p ValueParams) Ints(key string) ([]int, error) {
	var ii []int
	ss := p.vals[key]
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			return ii, errors.Wrapf(err, "failed to parse '%v' to int", s)
		}
		ii = append(ii, i)
	}
	return ii, nil
}

// Bool gets one value associated with key and returns it as bool
func (p ValueParams) Bool(key string) (bool, error) {
	v := p.vals.Get(key)
	vb, err := strconv.ParseBool(v)
	return vb, errors.Wrapf(err, "failed to parse '%v' to bool", v)
}

// Bools returns all value associated with key as bool slice
func (p ValueParams) Bools(key string) ([]bool, error) {
	var bb []bool
	ss := p.vals[key]
	for _, s := range ss {
		b, err := strconv.ParseBool(s)
		if err != nil {
			return bb, errors.Wrapf(err, "failed to parse '%v' to bool", s)
		}
		bb = append(bb, b)
	}
	return bb, nil
}
