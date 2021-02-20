package ctpm

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"os"
	"path/filepath"
)

type PublishOptions struct {
	Exclude []string
	Include []string
}

var PublishDefaultOptions = PublishOptions{
	Exclude: []string{".git/**", ".c3pm/**"},
	Include: []string{"c3pm.yml"},
}

func isFileInList(path string, rules []string) (bool, error) {
	for _, i := range rules {
		ok, err := doublestar.Match(i, path)
		if err != nil {
			return false, fmt.Errorf("failed to match [%s] with [%s] regex: %w", path, i, err)
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func getFilesFromRules(included []string, excluded []string) ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("publish list files walk: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		isIncluded, err := isFileInList(path, included)
		if err != nil {
			return fmt.Errorf("could not find if path [%s] is included: %w", path, err)
		}
		if isIncluded {
			isExcluded, err := isFileInList(path, excluded)
			if err != nil {
				return fmt.Errorf("could not find if path [%s] is excluded: %w", path, err)
			}
			if !isExcluded {
				files = append(files, path)
			}
		}
		return nil
	})
	return files, err
}

// Publish function makes an array of the files to include in the tarball
// based on the Include and Exclude fields of the c3pm.yaml
// The array is then given to the Upload function in the client
func Publish(pc *config.ProjectConfig, client api.API) error {
	included := append(pc.Manifest.Include, PublishDefaultOptions.Include...)
	included = append(included, pc.Manifest.Files.Sources...)
	included = append(included, pc.Manifest.Files.Includes...)
	excluded := append(pc.Manifest.Exclude, PublishDefaultOptions.Exclude...)

	files, err := getFilesFromRules(included, excluded)
	if err != nil {
		return fmt.Errorf("failed to list files to publish: %w", err)
	}
	return client.Upload(files)
}
