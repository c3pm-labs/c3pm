package ctpm

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/api"
	"github.com/c3pm-labs/c3pm/config"
	"io/ioutil"
	"os"
)

func Login(client api.API, login, password string) error {
	apiKey, err := client.Login(login, password)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(config.AuthFilePath(), []byte(apiKey), 0644)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to write auth file %s: %w", config.AuthFilePath(), err)
	}
	err = os.Mkdir(config.GlobalC3pmDirPath(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create .c3pm at %s: %w", config.GlobalC3pmDirPath(), err)
	}
	return ioutil.WriteFile(config.AuthFilePath(), []byte(apiKey), 0644)
}
