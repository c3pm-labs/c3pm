package cmakegen

import (
	"runtime"
	"strings"
)

func addPlatformSpecificCmake(base string, v CMakeVars) string {
	var tmpl strings.Builder
	tmpl.WriteString(base)
	if runtime.GOOS == "linux" && v.LinuxConfig != nil {
		addLinuxData(&tmpl, v)
	}
	return tmpl.String()
}
