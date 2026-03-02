package api

import (
	"net/http"
)

func NewRouter() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", HealthCheck)

	return mux
}
