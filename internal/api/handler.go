package api

import "net/http"


func HandleSSE(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Hello Get World"))

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

