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

func GenerateBuildFiles(sourceDir string, buildDir string, variables map[string]string) error {
	args := []string{
		"-S", sourceDir,
		"-B", buildDir,
	}
	for key, value := range variables {
		args = append(args, fmt.Sprintf("-D%s=%s", key, value))
	}
	return executeCmakeCLI(args...)
}

func Build(buildDir string) error {
	return executeCmakeCLI("--build", buildDir, "--config", "Release")
}

func Install(buildDir string) error {
	return executeCmakeCLI("--install", buildDir)
}
