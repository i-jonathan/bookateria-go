package main

import (
	"bookateria-api-go/document"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main()  {
	router := mux.NewRouter()
	// Routes for Documents
	subRouter := router.PathPrefix("/document").Subrouter()
	subRouter.HandleFunc("", document.GetDocuments).Methods("GET")
	subRouter.HandleFunc("/{id}", document.GetDocument).Methods("GET")
	subRouter.HandleFunc("", document.PostDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", document.UpdateDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", document.DeleteDocument).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":5000", router))
}
