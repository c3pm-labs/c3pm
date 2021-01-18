package input

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config/manifest"
)

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
	Name string
	Executable bool
	Library bool
	Description string
	Version string
	License string
}

func Init() (manifest.Manifest, error) {
	man := manifest.New()
	err := survey.Ask(InitSurvey, &man, SurveyOptions...)
	return man, err
}

func InitNonInteractive(val InitValues) (manifest.Manifest, error) {
	man := manifest.New()
	man.Name = val.Name
	if val.Executable {
		man.Type = manifest.Executable
	} else if val.Library {
		man.Type = manifest.Library
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
