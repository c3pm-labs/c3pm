package defaultadapter

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

{{- $c3pmGlobalDir:=.C3PMGlobalDir}}

target_sources({{.ProjectName}} PRIVATE {{.Sources}} {{.Headers}})
target_include_directories(
	{{- .ProjectName}} PRIVATE {{.IncludeDirs}}
	{{- range .Dependencies }}
		{{- $name:=.Name }}
		{{- $version:=.Version}}
		{{- range .IncludeDirs }}
			{{ $c3pmGlobalDir }}/cache/{{$name}}/{{$version}}/{{.}}
		{{- end }}
	{{- end }}
)
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

func library(v cmakeVars) (string, error) {
	funcMap := template.FuncMap{
		"AddTrailingSlash": func(text string) string {
			if !strings.HasSuffix(text, "/") {
				return text + "/"
			}
			return text
		},
	}
	cmake := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("cmakeLibrary").Funcs(funcMap).Parse(addPlatformSpecificCMake(libraryTemplate, v))
	if err != nil {
		return "", fmt.Errorf("could not parse cmake template: %w", err)
	}
	if err := tmpl.Execute(cmake, v); err != nil {
		return "", fmt.Errorf("could not template cmake: %w", err)
	}
	cmakeClean := cmake.String()
	if len(v.ExportedDir) == 0 {
		cmakeClean = removeCommand(cmakeClean, "install(DIRECTORY")
	}
	return cmakeClean, nil
}
