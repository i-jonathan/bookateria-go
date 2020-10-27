package forum

import "github.com/gorilla/mux"

func Router(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("/question").Subrouter()
	subRouter.HandleFunc("/all", GetQuestions).Methods("GET")
	subRouter.HandleFunc("/{id}", GetQuestion).Methods("GET")
	subRouter.HandleFunc("", PostQuestion).Methods("POST")
	subRouter.HandleFunc("/{id}", UpdateQuestion).Methods("POST")
	subRouter.HandleFunc("/{id}", DeleteQuestion).Methods("DELETE")
	subRouter.HandleFunc("/up-votes", GetQuestionUpVotes).Methods("GET")
	subRouter.HandleFunc("/up-votes", PostQuestionUpVote).Methods("POST")
	subRouter.HandleFunc("/up-votes/{id}", DeleteQuestionUpvote).Methods("DELETE")

	subRouter = router.PathPrefix("/answer").Subrouter()
	subRouter.HandleFunc("/all", GetAnswers).Methods("GET")
	subRouter.HandleFunc("/{id}", GetAnswer).Methods("GET")
	subRouter.HandleFunc("", PostAnswer).Methods("POST")
	subRouter.HandleFunc("/{id}", UpdateAnswer).Methods("POST")
	subRouter.HandleFunc("/{id}", DeleteAnswer).Methods("DELETE")
	subRouter.HandleFunc("/up-votes", GetAnswerUpVotes).Methods("GET")
	subRouter.HandleFunc("/up-votes", PostAnswerUpVote).Methods("POST")
	subRouter.HandleFunc("/up-votes/{id}", DeleteAnswerUpvote).Methods("DELETE")

	return router
}
