package ctpm

import (
	"archive/tar"
	"fmt"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/dependencies"
	"github.com/c3pm-labs/c3pm/env"
	"github.com/c3pm-labs/c3pm/registry"
	"io"
	"os"
	"path/filepath"
)

type InstallHandler struct{}

func (i InstallHandler) FetchDeps(request dependencies.PackageRequest) (dependencies.Dependencies, error) {
	libPath := config.LibCachePath(request.Name, request.Version)
	pc, err := config.Load(libPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	ret := make(dependencies.Dependencies)
	for k, v := range pc.Manifest.Dependencies {
		ret[k] = v
	}
	return ret, nil
}

func (i InstallHandler) PreAct(request dependencies.PackageRequest) error {
	fmt.Printf("Installing %s:%s\n", request.Name, request.Version)
	return Install(request.Name, request.Version)
}

func (i InstallHandler) PostAct(request dependencies.PackageRequest) error {
	fmt.Printf("Building %s:%s\n", request.Name, request.Version)
	// Build the lib
	libPath := config.LibCachePath(request.Name, request.Version)
	pc, err := config.Load(libPath)
	if err != nil {
		return fmt.Errorf("failed to read c3pm.yml: %w", err)
	}
	return Build(pc)
}

// Install fetches the package, unpacks it in the c3pm cache and builds it. If the lib already is
// in the cache, we don't do anything
func Install(name, version string) error {
	libPath := config.LibCachePath(name, version)
	// return early if lib is already in cache
	if _, err := os.Stat(libPath); !os.IsNotExist(err) {
		return nil
	}
	// fetch the tarball
	tarball, err := registry.FetchPackage(name, version, registry.Options{
		RegistryURL: env.REGISTRY_ENDPOINT,
	})
	if err != nil {
		return fmt.Errorf("error fetching package: %w", err)
	}
	// unpack the tarball
	err = unpackTarball(tarball, libPath)
	return err
}

func unpackTarball(pkg *os.File, dest string) error {
	tr := tar.NewReader(pkg)
	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating package directory: %w", err)
	}
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			_ = os.Remove(pkg.Name())
			return nil
		case err != nil:
			return fmt.Errorf("error unpacking tarball: %w", err)
		case header == nil:
			continue
		}
		target := filepath.Join(dest, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, os.ModePerm); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Dir(target), os.ModePerm)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			_ = f.Close()
		case tar.TypeSymlink:
			err = os.MkdirAll(filepath.Dir(target), os.ModePerm)
			if err != nil {
				return err
			}
			if err = os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		}
	}
}
