package manifest

import (
	"errors"
	"strconv"
	"strings"
)

//C3PMVersion holds the version of the c3pm manifest to use.
type C3PMVersion uint32

const (
	// C3PMVersion1 is the most recent manifest version.
	C3PMVersion1 = C3PMVersion(1)
)

// ErrParseC3PMVersion is returned when not able to parse the C3PM version to use.
var ErrParseC3PMVersion = errors.New("Could not parse C3PM version")

// C3PMVersionFromString returns the C3PMVersion corresponding to a given string.
func C3PMVersionFromString(version string) (C3PMVersion, error) {
	if !strings.HasPrefix(version, "v") {
		return C3PMVersion(0), ErrParseC3PMVersion
	}
	version = strings.TrimLeft(version, "v")
	val, err := strconv.ParseUint(version, 10, 32)
	if err != nil {
		return C3PMVersion(0), ErrParseC3PMVersion
	}
	return C3PMVersion(val), nil
}

// UnmarshalYAML is used to read the version from the manifest YAML file.
func (v *C3PMVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var version string
	err := unmarshal(&version)
	if err != nil {
		return err
	}
	*v, err = C3PMVersionFromString(version)
	if err != nil {
		return err
	}
	return nil
}

// String is used to return the stringified representation of the version.
func (v C3PMVersion) String() string {
	return "v" + strconv.FormatUint(uint64(v), 10)
}

// MarshalYAML is used to write the version as YAML.
func (v C3PMVersion) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}
