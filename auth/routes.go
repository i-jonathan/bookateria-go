package auth

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/refresh", RefreshToken).Methods("GET")
	router.HandleFunc("/logout", Logout).Methods("POST")
	//router.Use(AuthorizationMiddleware)
	router.HandleFunc("/login", SignIn).Methods("POST")
	return router
}
