package defaultadapter

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

var executableTemplate = `cmake_minimum_required(VERSION 3.16)
project({{.ProjectName}} VERSION {{.ProjectVersion}})

set(CMAKE_CXX_STANDARD {{.LanguageStandard}})
set(C3PM_PROJECT_NAME {{.ProjectName}})
set(C3PM_GLOBAL_DIR {{.C3PMGlobalDir}})

add_executable({{.ProjectName}})

target_sources({{.ProjectName}} PRIVATE
	{{.Sources}}
	{{.Headers}}
)
{{$c3pmGlobalDir:=.C3PMGlobalDir}}

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
{{range .Dependencies}}
find_library({{ .Name | ToUpper}} {{.Name}} "{{$c3pmGlobalDir}}/cache/{{.Name}}/{{.Version}}/")
{{end}}

{{.DependenciesConfig}}

target_link_libraries(
	{{.ProjectName}}
	PUBLIC
	{{range .Dependencies}}
	$<$<NOT:$<STREQUAL:"{{"${"}}{{.Name|ToUpper}}{{"}"}}","{{.Name|ToUpper}}-NOTFOUND">>:{{"${"}}{{.Name|ToUpper}}{{"}"}}>
	{{- end}}
)
`

func executable(v cmakeVars) (string, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
	}
	cmake := bytes.NewBuffer([]byte{})
	tmpl, err := template.New("cmakeExecutable").Funcs(funcMap).Parse(addPlatformSpecificCMake(executableTemplate, v))
	if err != nil {
		return "", fmt.Errorf("could not parse cmake template: %w", err)
	}
	if err := tmpl.Execute(cmake, v); err != nil {
		return "", fmt.Errorf("could not template cmake: %w", err)
	}
	return cmake.String(), nil
}
