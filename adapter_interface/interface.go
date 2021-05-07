package adapter_interface

import (
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
)

type Adapter interface {
	// Build builds the targets
	Build(pc *config.ProjectConfig) error
	// Targets return the paths of the targets built by the Build function
	Targets(pc *config.ProjectConfig) (targets []string, err error)
	//CmakeConfig return a string to add in the cmake of the project who use the library
	CmakeConfig(pc *config.ProjectConfig) (string, error)
}

type AdapterGetter interface {
	FromPC(adp *manifest.AdapterConfig) (Adapter, error)
}
