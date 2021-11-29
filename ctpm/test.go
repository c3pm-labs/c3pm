package ctpm

import (
	"errors"
	"fmt"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/adapter_interface"
	"github.com/c3pm-labs/c3pm/config"
)

func Test(pc *config.ProjectConfig) error {
	getter := adapter.AdapterGetterImp{}
	adp, err := getter.FromPC(pc.Manifest.Build.Adapter)
	if err != nil {
		return err
	}
	adpt, ok := adp.(adapter_interface.AdapterTestable)
	if !ok {
		return errors.New("Adapter does not support testing")
	}
	return adpt.Test(pc)
}

func AddDependenciesAndTest(pc *config.ProjectConfig) error {
	err := getAllDependencies(pc)
	if err != nil {
		return fmt.Errorf("error installing dependencies: %w", err)
	}
	err = Test(pc)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	return nil
}
