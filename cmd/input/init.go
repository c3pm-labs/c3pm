package input

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config/manifest"
)

//InitSurvey is the definition of the user interaction happening in the init command to choose project parameters.
var InitSurvey = []*survey.Question{
	{
		Name:     "Name",
		Prompt:   &survey.Input{Message: "Project name"},
		Validate: survey.Required,
	},
	{
		Name: "Type",
		Prompt: &survey.Select{Message: "Project type",
			Options: []string{
				manifest.Executable.String(),
				manifest.Library.String(),
			},
		},
	},
	{
		Name:   "Description",
		Prompt: &survey.Input{Message: "Project description"},
	},
	{
		Name: "Version",
		Prompt: &survey.Input{
			Message: "Project version",
			Default: "1.0.0",
		},
		Validate: func(ans interface{}) error {
			_, err := semver.NewVersion(ans.(string))
			return err
		},
		Transform: func(ans interface{}) (newAns interface{}) {
			v, _ := semver.NewVersion(ans.(string))
			return v.String()
		},
	},
	{
		Name: "License",
		Prompt: &survey.Input{
			Message: "Project license",
			Default: "UNLICENSED",
			Help:    "You can read about code licenses on https://choosealicense.com/",
		},
	},
}

type InitValues = struct {
	Name        string
	Type        string
	Description string
	Version     string
	License     string
}

//Init handles the user interaction happening during the init command
func Init() (manifest.Manifest, error) {
	man := manifest.New()
	err := survey.Ask(InitSurvey, &man, SurveyOptions...)
	return man, err
}

func InitNonInteractive(val InitValues) (manifest.Manifest, error) {
	var err error
	man := manifest.New()
	man.Name = val.Name
	man.Type, err = manifest.TypeFromString(val.Type)
	if err != nil {
		return manifest.Manifest{}, fmt.Errorf("failed to parse project type: %w", err)
	}
	man.Description = val.Description
	v, err := manifest.VersionFromString(val.Version)
	if err != nil {
		return man, fmt.Errorf("invalid version given: %w", err)
	}
	man.Version = v
	man.License = val.License
	return man, nil
}
