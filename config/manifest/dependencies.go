package manifest

import (
	"errors"
	"fmt"
	"strings"
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

func (d Dependencies) String() string {
	sb := strings.Builder{}
	i := 0
	for n, v := range d {
		if i != 0 {
			sb.WriteString(",")
			i = 1
		}
		sb.WriteString(fmt.Sprintf("[%s/%s]", n, v))
	}
	return sb.String()
}
