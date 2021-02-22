package ctpm

import (
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/config"
)

func Build(pc *config.ProjectConfig) error {
	adp, err := adapter.FromPc(pc)
	if err != nil {
		return err
	}
	return adp.Build(pc)
}
