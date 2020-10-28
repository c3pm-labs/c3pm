package manifest

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"reflect"
)

type Version struct {
	*semver.Version
}

func VersionFromString(version string) (v Version, err error) {
	v.Version, err = semver.NewVersion(version)
	if err != nil {
		return v, err
	}
	return v, nil
}

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

func (v Version) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}

func (v *Version) WriteAnswer(name string, value interface{}) (err error) {
	vs, ok := value.(string)
	if !ok {
		return fmt.Errorf("Version is not a string (" + reflect.TypeOf(vs).String() + ")")
	}
	v.Version, err = semver.NewVersion(vs)
	return err
}
