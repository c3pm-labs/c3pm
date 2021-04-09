package ctpm

import (
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"os"
	"path/filepath"
)

func isFileInList(path string, rules []string) (bool, error) {
	for _, rule := range rules {
		ok, err := doublestar.Match(rule, path)
		if err != nil {
			return false, fmt.Errorf("failed to match [%s] with [%s] regex: %w", path, rule, err)
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
// We enforce the exclusion of the .git and .c3pm directories and we enforce
// the inclusion of the c3pm.yml file
func Publish(pc *config.ProjectConfig, client api.API) error {
	fmt.Println("Manifest Include:", pc.Manifest.Publish.Include)
	included := append(pc.Manifest.Publish.Include, "c3pm.yml")
	excluded := append(pc.Manifest.Publish.Exclude, ".git/**", ".c3pm/**")

	fmt.Println("included:", included)
	fmt.Println("excluded:", excluded)
	files, err := getFilesFromRules(included, excluded)
	fmt.Println("files:", files)
	if err != nil {
		return fmt.Errorf("failed to list files to publish: %w", err)
	}
	return client.Upload(files)
}
