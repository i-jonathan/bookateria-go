package document

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("", GetDocuments).Methods("GET")
	router.HandleFunc("/{id}", GetDocument).Methods("GET")
	router.HandleFunc("", PostDocument).Methods("POST")
	router.HandleFunc("/{id}", UpdateDocument).Methods("POST")
	router.HandleFunc("/{id}", DeleteDocument).Methods("DELETE")

	return router
}