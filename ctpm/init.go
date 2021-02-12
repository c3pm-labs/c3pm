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
const readMeTemplate = `<p align="center">
<a href="https://c3pm.io/">
	<img alt="c3pm" src="https://dev.c3pm.io/assets/c3pm.png" width="546"/>
</a>
</p>

<p align="center">
	Your toolkit to dive into C++ easily
</p>

---

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [What is c3pm?](#what-is-c3pm)
- [Installing c3pm](#installing-c3pm)
- [Usage](#usage)
- [Start your project](#start-your-project)
- [Add a package](#add-a-package)
- [Building your project](#building-your-project)
- [Publishing your project](#publishing-your-project)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## What is c3pm?

**c3pm** stands for C++ package manager.

c3pm abstracts your build system and eases the management of your dependencies.

It manages your CMake and compiles your project with **minimal configuration**.

Feel free to explore the available packages on our [platform](https://c3pm.io).

## Installing c3pm

To install c3pm, click on this [link](https://github.com/c3pm-labs/c3pm/releases) and download the version
relase you want for your operating system.

c3pm is available for macOS, Windows and Linux.

## Usage

Once you installed c3pm, you can start using it.
Here are some basic commands you will need.

### Start your project
` + "```" + `shell
$ ctpm init
` + "```" + `

### Add a package
` + "```" + `shell
$ ctpm add <package>
` + "```" + `


### Building your project
` + "```" + `shell
$ ctpm build
` + "```" + `

### Publishing your project

` + "```" + `shell
$ ctpm publish
` + "```" + `

<br />

You can find a more complete list of the available commands [here](https://github.com/gabrielcolson/c3pm/tree/master/specs/cli) !
`

func generateReadMe(pc *config.ProjectConfig) error {
	return ioutil.WriteFile(filepath.Join(pc.ProjectRoot, "README.md"), []byte(readMeTemplate), 0644)
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
