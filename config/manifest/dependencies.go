package manifest

import (
	"errors"
)

type Dependencies map[string]string

var ErrParseDependencies = errors.New("Could not parse dependencies")

func DependenciesFromMap(dependencies map[string]string) (Dependencies, error) {
	if dependencies == nil {
		return Dependencies(nil), ErrParseDependencies
	}
	return Dependencies(dependencies), nil
}

func (d Dependencies) MarshalYAML() (interface{}, error) {
	return d, nil
}
