// The role of this package is to create simple mock servers to be used when testing interaction with C3PM's API.
package apitest

import "net/http"

func createServer(code int) http.Handler {
	sm := http.NewServeMux()
	sm.HandleFunc("/packages/publish", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
	sm.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
	return sm
}

//MockServer returns an HTTP server always responding with HTTP 200 response codes
func MockServer() http.Handler {
	return createServer(http.StatusOK)
}

//ErrorServer returns an HTTP server always responding with HTTP 507 response codes
func ErrorServer() http.Handler {
	return createServer(http.StatusInsufficientStorage)
}
