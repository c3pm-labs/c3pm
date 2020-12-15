package api_test

import (
	"fmt"
	"github.com/c3pm-labs/c3pm/api"
	"net/http"
)

func ExampleAPI_Login() {
	c := api.New(&http.Client{}, "xxxx")
	key, err := c.Login("login", "password")
	fmt.Println(key, err)
	// Output: auth_XXXX, nil
}

func ExampleAPI_Upload() {
	c := api.New(&http.Client{}, "xxxx")
	err := c.Upload([]string{"../../testhelpers/tars/single.tar"})
	fmt.Println(err)
	// Output: nil
}
