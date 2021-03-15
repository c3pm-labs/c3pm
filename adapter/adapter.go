package adapter

import (
	"errors"
	"github.com/c3pm-labs/c3pm/adapter/defaultadapter"
	"github.com/c3pm-labs/c3pm/config"
)

type Adapter interface {
	// Build builds the targets
	Build(pc *config.ProjectConfig) error
	// Targets return the paths of the targets built by the Build function
	Targets(pc *config.ProjectConfig) (targets []string, err error)
}

func FromPC(pc *config.ProjectConfig) (Adapter, error) {
	adp := pc.Manifest.Build.Adapter

	switch {
	case adp.Name == "c3pm" && adp.Version.String() == "0.0.1":
		return defaultadapter.New(), nil
	default:
		return nil, errors.New("only default adapter is supported")
	}
}
