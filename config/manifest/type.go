package manifest

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"reflect"
)

// Type holds the kind of project to manage.
type Type string

const (
	// Executable is the Type of a project yielding a binary executable.
	Executable = Type("executable")
	// Library is the Type of a project yielding a package to be included by others.
	Library = Type("library")
)

// String returns the string representation of a Type
func (t Type) String() string {
	return string(t)
}

// TypeFromString creates a Type instance from a corresponding string.
func TypeFromString(t string) (Type, error) {
	switch t {
	case Executable.String():
		return Executable, nil
	case Library.String():
		return Library, nil
	default:
		return Type(""), fmt.Errorf("invalid project type")
	}
}

// WriteAnswer allows type to be used in survey questions
func (t *Type) WriteAnswer(name string, value interface{}) (err error) {
	var ts string
	//nolint:gosimple
	switch value.(type) {
	case string:
		ts = value.(string)
	case survey.OptionAnswer:
		ts = value.(survey.OptionAnswer).Value
	default:
		return fmt.Errorf("Type is not a string (" + reflect.TypeOf(value).String() + ")")
	}
	*t, err = TypeFromString(ts)
	return err
}
