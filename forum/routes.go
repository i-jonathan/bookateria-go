package forum

import "github.com/gorilla/mux"

// Router - All routes for forum feature
func Router(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("/question").Subrouter()
	subRouter.HandleFunc("/all", QuestionSearch).Queries("search", "{search}").Methods("GET")
	subRouter.HandleFunc("/all", FilterQuestionByTags).Queries("filter", "{filter}").Methods("GET")
	subRouter.HandleFunc("/all", GetQuestions).Methods("GET")
	subRouter.HandleFunc("/{slug}", GetQuestion).Methods("GET")
	subRouter.HandleFunc("", PostQuestion).Methods("POST")
	subRouter.HandleFunc("/{slug}", UpdateQuestion).Methods("PUT")
	subRouter.HandleFunc("/{slug}", DeleteQuestion).Methods("DELETE")
	subRouter.HandleFunc("{slug}/up-votes", GetQuestionUpVotes).Methods("GET")
	subRouter.HandleFunc("{slug}/up-votes", PostQuestionUpVote).Methods("POST")
	subRouter.HandleFunc("{slug}/up-votes/{id}", DeleteQuestionUpvote).Methods("DELETE")

	subRouter = router.PathPrefix("/answer").Subrouter()
	subRouter.HandleFunc("/all", GetAnswers).Methods("GET")
	subRouter.HandleFunc("/{slug}", GetAnswer).Methods("GET")
	subRouter.HandleFunc("", PostAnswer).Methods("POST")
	subRouter.HandleFunc("/{slug}", UpdateAnswer).Methods("PUT")
	subRouter.HandleFunc("/{slug}", DeleteAnswer).Methods("DELETE")
	subRouter.HandleFunc("{slug}/up-votes", GetAnswerUpVotes).Methods("GET")
	subRouter.HandleFunc("{slug}/up-votes", PostAnswerUpVote).Methods("POST")
	subRouter.HandleFunc("{slug}/up-votes/{id}", DeleteAnswerUpvote).Methods("DELETE")
	return router
}
