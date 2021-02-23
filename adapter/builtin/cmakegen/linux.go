package cmakegen

import (
	"strings"
)

var pthreadTemplate = `
set(THREADS_PREFER_PTHREAD_FLAG ON)
find_package(Threads REQUIRED)
target_link_libraries({{.ProjectName}} PUBLIC Threads::Threads)
`

func addLinuxData(sb *strings.Builder, v CMakeVars) {
	if v.LinuxConfig.UsePthread {
		sb.WriteString(pthreadTemplate)
	}
}
