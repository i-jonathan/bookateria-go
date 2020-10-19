package account

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/all", AllUsers).Methods("GET")
	router.HandleFunc("", PostUser).Methods("POST")
	router.HandleFunc("/{id}", GetUser).Methods("GET")

	return router
}