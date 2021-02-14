package manifest

import (
	"errors"
)

// Dependencies holds the list of the dependencies of a project.
type Dependencies map[string]string

// ErrParseDependencies is returned when an error occured while reading dependencies.
var ErrParseDependencies = errors.New("Could not parse dependencies")

// TODO: Unused
func DependenciesFromMap(dependencies map[string]string) (Dependencies, error) {
	if dependencies == nil {
		return Dependencies(nil), ErrParseDependencies
	}
	return Dependencies(dependencies), nil
}

// MarshalYAML is used to write dependencies as YAML.
func (d Dependencies) MarshalYAML() (interface{}, error) {
	return d, nil
}
