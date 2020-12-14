package apitest

import "net/http"

func MockServer() http.Handler {
	sm := http.NewServeMux()
	sm.HandleFunc("/auth/publish", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	})
	sm.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
	})
	return sm
}

func ErrorServer() http.Handler {
	sm := http.NewServeMux()
	sm.HandleFunc("/auth/publish", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInsufficientStorage)
	})
	sm.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInsufficientStorage)
	})
	return sm
}
