package assignment

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/all", GetQuestions).Methods("GET")
	router.HandleFunc("/add", PostQuestion).Methods("POST")
	router.HandleFunc("/{slug}", GetQuestion).Methods("GET")
	router.HandleFunc("/{slug}", UpdateQuestion).Methods("PUT")
	router.HandleFunc("/{slug}/delete", DeleteQuestion).Methods("DELETE")
	router.HandleFunc("/{qSlug}/submit", PostSubmission).Methods("POST")
	router.HandleFunc("/{qSlug}/submissions", GetSubmissions).Methods("GET")
	router.HandleFunc("/{qSlug}/submission/{aSlug}", GetSubmission).Methods("GET")
	return router
}
