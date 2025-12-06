// Package user
package user

import "net/http"

func NewRouter(storage *Storage) http.Handler {
	mux := http.NewServeMux()
	handler := NewHandler(storage)

	mux.HandleFunc("/users", handler.List)
	mux.HandleFunc("/users/create", handler.Create)

	return mux
}
