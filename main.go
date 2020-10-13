package main

import (
	"bookateria-api-go/documents"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main()  {
	router := mux.NewRouter()
	// Routes for Documents
	subRouter := router.PathPrefix("/document").Subrouter()
	subRouter.HandleFunc("", documents.GetDocuments).Methods("GET")
	subRouter.HandleFunc("/{id}", documents.GetDocument).Methods("GET")
	subRouter.HandleFunc("", documents.PostDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", documents.UpdateDocument).Methods("POST")
	subRouter.HandleFunc("/{id}", documents.DeleteDocument).Methods("DELETE")


	log.Fatal(http.ListenAndServe(":5000", router))
}

