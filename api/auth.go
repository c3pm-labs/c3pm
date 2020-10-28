package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (c *API) Login(login, password string) (string, error) {
	body, err := json.Marshal(struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	var data = struct {
		ApiKey string `json:"apiKey"`
	}{}
	err = c.fetch(http.MethodPost, "/auth/login", bytes.NewReader(body), &data)
	return data.ApiKey, err
}
