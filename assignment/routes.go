package assignment

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	router.HandleFunc("/questions", GetQuestions).Methods("GET")
	router.HandleFunc("/question/add", PostQuestion).Methods("POST")
	router.HandleFunc("/question/{slug}", GetQuestion).Methods("GET")
	router.HandleFunc("/question/{qSlug}/submit", PostSubmission).Methods("POST")
	return router
}
