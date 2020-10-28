package cmakegen

import (
	"bytes"
	"fmt"
	"text/template"
)

var executableTemplate = `cmake_minimum_required(VERSION 3.16)
project({{.ProjectName}} VERSION {{.ProjectVersion}})

set(CMAKE_CXX_STANDARD {{.LanguageStandard}})

add_executable({{.ProjectName}})

target_sources({{.ProjectName}} PRIVATE {{.Sources}} {{.Includes}})
{{$home:=.Home}}

target_include_directories(
	{{.ProjectName}} PRIVATE {{.IncludeDirs}}
	{{ range .Dependencies }}
		{{$name:=.Name}}
		{{$version:=.Version}}
		{{ range .ExportedIncludeDirs }}
			{{$home}}/.c3pm/cache/{{$name}}/{{$version}}/{{.}}
		{{ end }}
	{{end}}
)

target_link_libraries(
	{{.ProjectName}}
	PUBLIC
	{{ range .Dependencies }}
		-L{{$home}}/.c3pm/cache/{{.Name}}/{{.Version}}/lib
		{{ range .Targets }} -l{{.}} {{end}}
	{{end}}
)
`

func executable(v CMakeVars) (string, error) {
	cmake := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("cmakeExecutable").Parse(addPlatformSpecificCmake(executableTemplate, v))
	if err != nil {
		return "", fmt.Errorf("could not parse cmake template: %w", err)
	}
	if err := tmpl.Execute(cmake, v); err != nil {
		return "", fmt.Errorf("could not template cmake: %w", err)
	}
	return cmake.String(), nil
}
