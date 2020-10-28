package manifest

type CustomCmake struct {
	Path      string            `yaml:"path"`
	Variables map[string]string `yaml:"variables"`
	Targets   []string          `yaml:"targets"`
}
