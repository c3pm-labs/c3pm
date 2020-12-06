package ctpm

import (
	"github.com/c3pm-labs/c3pm/config"
	"os"
)

func Logout() error {
	return os.Remove(config.AuthFilePath())
}
