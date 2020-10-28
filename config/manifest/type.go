package manifest

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"reflect"
)

type Type string

const (
	Executable = Type("executable")
	Library    = Type("library")
)

func (t Type) String() string {
	return string(t)
}

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
