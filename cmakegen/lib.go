package cmakegen

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var libraryTemplate = `cmake_minimum_required(VERSION 3.16)
project({{.ProjectName}} VERSION {{.ProjectVersion}})

set(CMAKE_CXX_STANDARD {{.LanguageStandard}})

add_library({{.ProjectName}} STATIC)

target_sources({{.ProjectName}} PRIVATE {{.Sources}} {{.Includes}})
target_include_directories({{.ProjectName}} PRIVATE {{.IncludeDirs}})

install(
	DIRECTORY {{ .ExportedDir }}
	DESTINATION ${CMAKE_INSTALL_PREFIX}
)
install(TARGETS {{.ProjectName}} ARCHIVE)
`

func removeCommand(cmake string, command string) string {
	var cmakeClean string
	scanner := bufio.NewScanner(strings.NewReader(cmake))
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), command) {
			cmakeClean += scanner.Text() + "\n"
		}
	}
	return cmakeClean
}

func library(v CMakeVars) (string, error) {
	cmake := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("cmakeExecutable").Parse(addPlatformSpecificCmake(libraryTemplate, v))
	if err != nil {
		return "", fmt.Errorf("could not parse cmake template: %w", err)
	}
	if err := tmpl.Execute(cmake, v); err != nil {
		return "", fmt.Errorf("could not template cmake: %w", err)
	}
	cmakeClean := cmake.String()
	if len(v.IncludeDirs) == 0 {
		cmakeClean = removeCommand(cmakeClean, "target_include_directories")
	}
	if len(v.ExportedDir) == 0 {
		cmakeClean = removeCommand(cmakeClean, "install(DIRECTORY")
	}
	return cmakeClean, nil
}
