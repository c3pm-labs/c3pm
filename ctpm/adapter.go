package ctpm

import (
	"errors"
	"github.com/c3pm-labs/c3pm/adapter"
	"github.com/c3pm-labs/c3pm/adapter/builtin"
	"github.com/c3pm-labs/c3pm/config"
)

func adapterFromPc(pc *config.ProjectConfig) (adapter.Adapter, error) {
	if pc.Manifest.Build.Adapter == nil {
		return builtin.New(), nil
	} else {
		return nil, errors.New("only default adapter is supported")
	}
}
