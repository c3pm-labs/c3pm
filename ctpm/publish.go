package ctpm

import (
	"bufio"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"github.com/c3pm-labs/c3pm/config/manifest"
	"os"
	"path/filepath"
	"strings"
)

var IgnoreFiles = []string{".gitignore"}

type PublishOptions struct {
	Exclude []string
	Include []string
}

var PublishDefaultOptions = PublishOptions{
	Exclude: []string{".git/**", ".c3pm/**"},
	Include: []string{"c3pm.yml"},
}

func publishFiles(userList []string, filesType string, filesConfig manifest.FilesConfig) ([]string, error) {
	var files []string
	//var err error
	for _, ignoreFile := range IgnoreFiles {
		f, err := os.Open(ignoreFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to open ignore file [%s]: %w", ignoreFile, err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			files = append(files, scanner.Text())
		}
		f.Close()
	}
	files = append(files, userList...)
	var defaultOptions []string
	if filesType == "include" {
		defaultOptions = append(PublishDefaultOptions.Include, filesConfig.Sources...)
		defaultOptions = append(defaultOptions, filesConfig.Includes...)
	} else {
		defaultOptions = PublishDefaultOptions.Exclude
	}
	files = append(files, defaultOptions...)
	return files, nil
}

// IsFileInList return true if path is in rules, it will use rules/regex given
// rules follow pattern format from gitignore
func isFileInList(rules []string, path string) (bool, error) {
	var negate bool
	var fileIsInRules = false
	for _, i := range rules {
		if strings.HasPrefix(i, "!") {
			i = i[1:]
			negate = true
		} else {
			negate = false
		}
		ok, err := doublestar.Match(i, path)
		if err != nil {
			return false, fmt.Errorf("failed to match [%s] with [%s] regex: %w", path, i, err)
		}
		if ok {
			fileIsInRules = true
		}
		if ok && negate {
			fileIsInRules = false
		}
	}
	return fileIsInRules, nil
}

func publishListFiles(excluded []string, included []string) ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("publish list files walk: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		isIncluded, err := isFileInList(included, path)
		if err != nil {
			return fmt.Errorf("could not find if path [%s] is included: %w", path, err)
		}
		if !isIncluded {
			return nil
		}
		isExcluded, err := isFileInList(excluded, path)
		if err != nil {
			return fmt.Errorf("could not find if path [%s] is excluded: %w", path, err)
		}
		if isExcluded {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func Publish(pc *config.ProjectConfig, client api.API, opt PublishOptions) error {
	excluded, err := publishFiles(opt.Exclude, "exclude", pc.Manifest.Files)
	if err != nil {
		return fmt.Errorf("failed to parse excluded files: %w", err)
	}
	included, err := publishFiles(opt.Include, "include", pc.Manifest.Files)
	if err != nil {
		return fmt.Errorf("failed to parse included files: %w", err)
	}
	files, err := publishListFiles(excluded, included)
	if err != nil {
		return fmt.Errorf("failed to list files to publish: %w", err)
	}
	return client.Upload(files)
}
