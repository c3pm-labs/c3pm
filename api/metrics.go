package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *API) CountDownload(packageName string) error {
	body, err := json.Marshal(struct {
		PackageName string `json:"packageName"`
	}{
		PackageName: packageName,
	})
	if err != nil {
		return err
	}
	err = c.fetch(http.MethodPost, "/auth/login", bytes.NewReader(body), nil)
	if err != nil {
		return fmt.Errorf("failed to count download: %w", err)
	}
	return nil
}
