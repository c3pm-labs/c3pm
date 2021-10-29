package input

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"github.com/mitchellh/go-spdx"
	"strings"
)

var initLicenses []string = nil

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
			Suggest: func(toComplete string) []string {
				if initLicenses == nil {
					var err error
					initLicenses, err = listLicences()
					if err != nil {
						return nil
					}
				}
				filteredLicences := []string{}
				for _, l := range initLicenses {
					if strings.HasPrefix(l, toComplete) {
						filteredLicences = append(filteredLicences, l)
					}
				}
				return filteredLicences
			},
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

func listLicences() ([]string, error) {
	licenseList, err := spdx.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	licenses := []string{}
	for _, li := range licenseList.Licenses {
		licenses = append(licenses, li.ID)
	}
	return licenses, nil
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
