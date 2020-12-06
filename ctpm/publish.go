package ctpm

import (
	"bufio"
	"fmt"
	"github.com/bmatcuk/doublestar"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"os"
	"path/filepath"
	"strings"
)

var IgnoreFiles = []string{".gitignore", ".c3pmignore"}

type PublishOptions struct {
	Ignore []string
}

var PublishDefaultOptions = PublishOptions{
	Ignore: []string{".git/**", ".c3pm/**"},
}

func publishIgnoredFiles(opt PublishOptions) ([]string, error) {
	var ignored []string
	//var err error
	for _, ignoreFile := range IgnoreFiles {
		f, err := os.Open(ignoreFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to open ignore file [%s]: %w", ignoreFile, err)
		}
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			ignored = append(ignored, scanner.Text())
		}
		f.Close()
	}
	ignored = append(ignored, opt.Ignore...)
	ignored = append(ignored, PublishDefaultOptions.Ignore...)
	return ignored, nil
}

// isIgnored return true if path is ignored, it will use rules/regex given
// rules follow pattern format from gitignore
func isIgnored(rules []string, path string) (bool, error) {
	var negate bool
	var ignored = false
	for _, i := range rules {
		if strings.HasPrefix(i, "!") {
			i = i[1:]
			negate = true
		} else {
			negate = false
		}
		ok, err := doublestar.Match(i, path)
		if err != nil {
			return false, fmt.Errorf("failed to match [%s] with [%s] ignore regex: %w", path, i, err)
		}
		if ok {
			ignored = true
		}
		if ok && negate {
			ignored = false
		}
	}
	return ignored, nil
}

func publishListFiles(ignored []string) ([]string, error) {
	var files []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("publish list files walk: %w", err)
		}
		if info.IsDir() {
			return nil
		}
		ok, err := isIgnored(ignored, path)
		if err != nil {
			return fmt.Errorf("could not find if path [%s] is ignored: %w", path, err)
		}
		if ok {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func Publish(pc *config.ProjectConfig, client api.API, opt PublishOptions) error {
	ignored, err := publishIgnoredFiles(opt)
	if err != nil {
		return fmt.Errorf("failed to parse ignored files: %w", err)
	}
	files, err := publishListFiles(ignored)
	if err != nil {
		return fmt.Errorf("failed to list files to publish: %w", err)
	}
	return client.Upload(files)
}
