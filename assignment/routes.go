package assignment

import "github.com/gorilla/mux"

// Router is all assignment portal routes
// 		return *mux.Router
func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/all", getQuestions).Methods("GET")
	router.HandleFunc("/add", postQuestion).Methods("POST")
	router.HandleFunc("/{slug}", getQuestion).Methods("GET")
	router.HandleFunc("/{slug}", updateQuestion).Methods("PUT")
	router.HandleFunc("/{slug}/delete", deleteQuestion).Methods("DELETE")
	router.HandleFunc("/{qSlug}/submit", PostSubmission).Methods("POST")
	router.HandleFunc("/{qSlug}/submissions", getSubmissions).Methods("GET")
	router.HandleFunc("/{qSlug}/submission/{aSlug}", getSubmission).Methods("GET")
	return router
}
