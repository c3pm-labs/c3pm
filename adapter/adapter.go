package adapter

import (
	"github.com/c3pm-labs/c3pm/config"
)

type Adapter interface {
	// Build builds the targets
	Build(pc *config.ProjectConfig) error
	// Targets return the paths of the targets built by the Build function
	Targets(pc *config.ProjectConfig) (targets []string, err error)
}
