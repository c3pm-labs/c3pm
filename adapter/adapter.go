package adapter

import (
	"errors"
	"github.com/c3pm-labs/c3pm/adapter/cmake_adapter"
	"github.com/c3pm-labs/c3pm/adapter/defaultadapter"
	"github.com/c3pm-labs/c3pm/adapter/irrlichtadapter"
	"github.com/c3pm-labs/c3pm/adapter/libuv_adapter"
	"github.com/c3pm-labs/c3pm/adapter_interface"
	"github.com/c3pm-labs/c3pm/config/manifest"
)

type AdapterGetterImp struct{}

func (AdapterGetterImp) FromPC(adp *manifest.AdapterConfig) (adapter_interface.Adapter, error) {
	switch {
	case adp.Name == "c3pm" && adp.Version.String() == "0.0.1":
		return defaultadapter.New(AdapterGetterImp{}), nil
	case adp.Name == "irrlicht" && adp.Version.String() == "0.0.1":
		return irrlichtadapter.New(), nil
	case adp.Name == "cmake" && adp.Version.String() == "0.0.1":
		return cmake_adapter.New(), nil
	case adp.Name == "uv" && adp.Version.String() == "0.0.1":
		return libuv_adapter.New(), nil
	default:
		return nil, errors.New("only default adapter is supported")
	}
}
