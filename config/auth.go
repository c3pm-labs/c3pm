package config

import (
	"fmt"
	"io/ioutil"
	"path"
)

func TokenStrict() (string, error) {
	content, err := ioutil.ReadFile(AuthFilePath())
	if err != nil {
		return "", fmt.Errorf("failed to read auth file: %w", err)
	}
	return string(content), nil
}

// Token returns an empty string if it fails to read token
func Token() string {
	token, _ := TokenStrict()
	return token
}

func AuthFilePath() string {
	return path.Join(GlobalC3pmDirPath(), "auth.cfg")
}
