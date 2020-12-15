// The goal of the cmake package is to handle interaction with the CMake Command Line Interface.
// CMake is used internally by C3PM to manage the build and installation phases of using a C3PM project.
//
// More information about what the CMake CLI does can be found on CMake's website: https://cmake.org/cmake/help/latest/manual/cmake.1.html
package cmake

import (
	"fmt"
	"os"
	"os/exec"
)

func executeCmakeCLI(args ...string) error {
	cmd := exec.Command("cmake", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start cmake: %w", err)
	}
	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("cmake process failed: %w", err)
	}
	return nil
}

//GenerateBuildFiles runs the cmake CLI to generate CMake build files.
//C3PM uses CMake's -S option for setting the source directory, the -B option for the build directory, and the -D option for setting build variables.
//
//See CMake's documentation for more information: https://cmake.org/cmake/help/latest/manual/cmake.1.html#generate-a-project-buildsystem
func GenerateBuildFiles(sourceDir, buildDir string, variables map[string]string) error {
	args := []string{
		"-S", sourceDir,
		"-B", buildDir,
	}
	for key, value := range variables {
		args = append(args, fmt.Sprintf("-D%s=%s", key, value))
	}
	return executeCmakeCLI(args...)
}

//Build runs the CMake CLI to build a C3PM project
//
//See CMake's documentation for more information: https://cmake.org/cmake/help/latest/manual/cmake.1.html#build-a-project
func Build(buildDir string) error {
	return executeCmakeCLI("--build", buildDir, "--config", "Release")
}

//Install runs the CMake CLI to install a project.
//
//See CMake's documentation for more information: https://cmake.org/cmake/help/latest/manual/cmake.1.html#install-a-project
func Install(buildDir string) error {
	return executeCmakeCLI("--install", buildDir)
}
