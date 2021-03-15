package defaultadapter

import (
	"runtime"
	"strings"
)

func addPlatformSpecificCMake(base string, v cmakeVars) string {
	var tmpl strings.Builder
	tmpl.WriteString(base)
	if runtime.GOOS == "linux" && v.LinuxConfig != nil {
		addLinuxData(&tmpl, v)
	}
	return tmpl.String()
}
