package ctpm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/c3pm-labs/c3pm/cmakegen"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/mitchellh/go-spdx"
)

type InitOptions struct {
	NoTemplate bool
}

var InitDefaultOptions = InitOptions{
	NoTemplate: false,
}

func Init(pc *config.ProjectConfig, opt InitOptions) error {
	err := os.MkdirAll(pc.ProjectRoot, os.ModePerm) // umask will take care of permissions
	if err != nil {
		return err
	}
	if pc.Manifest.License != "UNLICENSED" {
		err = generateLicenseFile(pc)
		if err != nil {
			return err
		}
	}
	if !opt.NoTemplate {
		err = generateReadMe(pc)
		if err != nil {
			return err
		}
		if pc.Manifest.Type == "executable" {
			err := saveExecutableTemplate(pc)
			if err != nil {
				return err
			}
		} else {
			err := saveLibraryTemplate(pc)
			if err != nil {
				return err
			}
		}
	}
	err = pc.Save()
	if err != nil {
		return fmt.Errorf("failed to save project file: %w", err)
	}
	return cmakegen.Generate(pc)
}

const execTemplate = `#include <iostream>

int main() {
	std::cout << "Hello c3pm!" << std::endl;
	return 0;
}`
const includeTemplate = `#pragma once

void hello();`
const libTemplate = `#include <iostream>
#include "{{.Name}}.hpp"

void hello() {
	std::cout << "Hello c3pm!" << std::endl;
}`
const readMeTemplate = `

A new C++ project.

## Getting Started

This project is a starting point for a C++ project.

A few helpful commands to get you started if this is your first time using c3pm:

### Building your project
` + "```" + `shell
$ ctpm build
` + "```" + `

### Add a package
` + "```" + `shell
$ ctpm add <package>
` + "```" + `

### Publishing your project
` + "```" + `shell
$ ctpm publish
` + "```" + `

<br />

For help getting started with c3pm, view our
[online documentation](https://docs.c3pm.io/), which offers tutorials, samples and
a list of all available commands.
`

func generateReadMe(pc *config.ProjectConfig) error {
	projectName := `# ` + pc.Manifest.Name
	return ioutil.WriteFile(filepath.Join(pc.ProjectRoot, "README.md"), []byte(projectName+readMeTemplate), 0644)
}

func generateLicenseFile(pc *config.ProjectConfig) error {
	lic, err := spdx.License(pc.Manifest.License)
	if err != nil {
		return err
	}
	if len(lic.Text) == 0 {
		return fmt.Errorf("generate %s", lic.Name)
	}
	return ioutil.WriteFile(filepath.Join(pc.ProjectRoot, "LICENSE"), []byte(lic.Text), 0644)
}

func saveExecutableTemplate(pc *config.ProjectConfig) error {
	if err := os.Mkdir(filepath.Join(pc.ProjectRoot, "src"), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(pc.ProjectRoot, "src", "main.cpp"), []byte(execTemplate), 0644)
}

func saveLibraryTemplate(pc *config.ProjectConfig) error {
	pc.Manifest.Files.IncludeDirs = append(pc.Manifest.Files.IncludeDirs, "include")
	t := template.Must(template.New("libTemplate").Parse(libTemplate))
	if err := os.Mkdir(filepath.Join(pc.ProjectRoot, "src"), os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(pc.ProjectRoot, "src", pc.Manifest.Name+".cpp"))
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.ExecuteTemplate(f, "libTemplate", pc.Manifest); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(pc.ProjectRoot, "include"), os.ModePerm); err != nil {
		return err
	}
	f, err = os.Create(filepath.Join(pc.ProjectRoot, "include", pc.Manifest.Name+".hpp"))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(includeTemplate))
	return err
}
