package ctpm

import (
	"github.com/gabrielcolson/c3pm/cli/config"
	"os"
)

func Logout() error {
	return os.Remove(config.AuthFilePath())
}
