package account

import "github.com/gorilla/mux"

// Router contains all endpoints for accounts.
func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/all", allUsers).Methods("GET")
	router.HandleFunc("", postUser).Methods("POST")
	router.HandleFunc("/{id}", getUser).Methods("GET")
	router.HandleFunc("/verify-email", verifyEmail).Methods("POST")
	router.HandleFunc("/request-otp", requestOTP).Methods("POST")

	return router
}
