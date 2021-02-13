package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Login handles the login logic with the API.
// It takes the user's email and password as strings,
// sends them to the API, and returns the ApiKey object
// to be used in further calls to the API.
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
