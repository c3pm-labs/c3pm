package manifest

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"reflect"
)

// Version holds the Semantic Versioning-compatible version of a package.
type Version struct {
	*semver.Version
}

// VersionFromString converts a string to a Version instance.
func VersionFromString(version string) (v Version, err error) {
	v.Version, err = semver.NewVersion(version)
	if err != nil {
		return v, err
	}
	return v, nil
}

// UnmarshalYAML is used to read a Version from a YAML file.
func (v *Version) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var version string
	err := unmarshal(&version)
	if err != nil {
		return err
	}
	v.Version, err = semver.NewVersion(version)
	if err != nil {
		return err
	}
	return nil
}

// MarshalYAML is used to write a Version to a YAML file.
func (v Version) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

// WriteAnswer allows type to be used in survey questions
func (v *Version) WriteAnswer(name string, value interface{}) (err error) {
	vs, ok := value.(string)
	if !ok {
		return fmt.Errorf("Version is not a string (" + reflect.TypeOf(vs).String() + ")")
	}
	v.Version, err = semver.NewVersion(vs)
	return err
}
