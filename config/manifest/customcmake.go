package manifest

// CustomCMake is used in cases where a project needs a complex custom CMake system.
type CustomCMake struct {
	Path      string            `yaml:"path"`
	Variables map[string]string `yaml:"variables"`
	Targets   []string          `yaml:"targets"`
}
