package builtin

import (
	"runtime"
	"strings"
)

func addPlatformSpecificCMake(base string, v CMakeVars) string {
	var tmpl strings.Builder
	tmpl.WriteString(base)
	if runtime.GOOS == "linux" && v.LinuxConfig != nil {
		addLinuxData(&tmpl, v)
	}
	return tmpl.String()
}
