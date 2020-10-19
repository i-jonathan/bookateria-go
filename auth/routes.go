package auth

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/login", SignIn).Methods("POST")
	return router
}