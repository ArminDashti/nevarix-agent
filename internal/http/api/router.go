package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(token string) http.Handler {
	r := mux.NewRouter()

	r.Use(authMiddleware(token))

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/config", getConfig).Methods("POST")
	v1.HandleFunc("/health", health).Methods("GET")
	v1.HandleFunc("/stat", stat).Methods("GET")

	return r
}
