package adapter

import (
	"errors"
	"github.com/c3pm-labs/c3pm/adapter/builtin"
	"github.com/c3pm-labs/c3pm/config"
)

type Adapter interface {
	// Build builds the targets
	Build(pc *config.ProjectConfig) error
	// Targets return the paths of the targets built by the Build function
	Targets(pc *config.ProjectConfig) (targets []string, err error)
}

func FromPc(pc *config.ProjectConfig) (Adapter, error) {
	if pc.Manifest.Build.Adapter == nil {
		return builtin.New(), nil
	} else {
		return nil, errors.New("only default adapter is supported")
	}
}
