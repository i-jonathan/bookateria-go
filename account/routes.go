package account

import "github.com/gorilla/mux"

// Router contains all endpoints for accounts.
func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/all", AllUsers).Methods("GET")
	router.HandleFunc("", PostUser).Methods("POST")
	router.HandleFunc("/{id}", GetUser).Methods("GET")
	router.HandleFunc("/verify-email", VerifyEmail).Methods("POST")
	router.HandleFunc("/request-otp", RequestOTP).Methods("POST")

	router.HandleFunc("/{id}", DeleteUser).Methods("DELETE")

	return router
}
