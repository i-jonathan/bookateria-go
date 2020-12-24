package auth

import "github.com/gorilla/mux"

// Router contains all routes for authorization
func Router(router *mux.Router) *mux.Router {
	// router.HandleFunc("/refresh", RefreshToken).Methods("GET")
	router.HandleFunc("/logout", Logout).Methods("POST")
	//router.Use(AuthorizationMiddleware)
	router.HandleFunc("/login", SignIn).Methods("POST")
	return router
}
