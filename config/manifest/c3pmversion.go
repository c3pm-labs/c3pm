package manifest

import (
	"errors"
	"strconv"
	"strings"
)

type C3pmVersion uint32

const (
	C3pmVersion1 = C3pmVersion(1)
)

var ErrParseC3pmVersion = errors.New("Could not parse C3PM version")

func C3pmVersionFromString(version string) (C3pmVersion, error) {
	if !strings.HasPrefix(version, "v") {
		return C3pmVersion(0), ErrParseC3pmVersion
	}
	version = strings.TrimLeft(version, "v")
	val, err := strconv.ParseUint(version, 10, 32)
	if err != nil {
		return C3pmVersion(0), ErrParseC3pmVersion
	}
	return C3pmVersion(val), nil
}

func (v *C3pmVersion) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var version string
	err := unmarshal(&version)
	if err != nil {
		return err
	}
	*v, err = C3pmVersionFromString(version)
	if err != nil {
		return err
	}
	return nil
}

func (v C3pmVersion) String() string {
	return "v" + strconv.FormatUint(uint64(v), 10)
}

func (v C3pmVersion) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}
