package handlers

import "net/http"

func RootHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Running API v1\n"))
	return
}
