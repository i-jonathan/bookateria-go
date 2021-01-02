package document

import "github.com/gorilla/mux"

// Router contains all routes for documents feature
func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("", SearchDocuments).Queries("search", "{search}").Methods("GET")
	router.HandleFunc("", FilterByTags).Queries("filter", "{filter}").Methods("GET")
	router.HandleFunc("", GetDocuments).Methods("GET")
	router.HandleFunc("/{id}", GetDocument).Methods("GET")
	router.HandleFunc("", PostDocument).Methods("POST")
	router.HandleFunc("/{id}", UpdateDocument).Methods("PUT")
	router.HandleFunc("/{id}", DeleteDocument).Methods("DELETE")

	return router
}
