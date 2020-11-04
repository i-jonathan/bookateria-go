package forum

import (
	"bookateria-api-go/account"
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	answerUpVotes   []AnswerUpvote
	answerUpVote    AnswerUpvote
	answers         []Answer
	answer          Answer
	db              = InitDatabase()
	question        Question
	questions       []Question
	questionUpVote  QuestionUpVote
	questionUpVotes []QuestionUpVote
	user            account.User
)

type Response struct {
	Message string `json:"message"`
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slugToFind, _ := params["slug"]

	if !XExists(slugToFind, "question") {
		// Checks if question exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Question Not Found"})
		log.Handler("info", "Question not found", err)
		return
	}

	db.Preload(clause.Associations).First(&question, slugToFind)
	err := json.NewEncoder(w).Encode(question)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func GetQuestions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Preload(clause.Associations).Find(&questions)
	err := json.NewEncoder(w).Encode(questions)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("warning", "JSON decoder error", err)
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}
	db.Find(&user, "email = ?", strings.ToLower(email))
	question.User = user
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	question.Slug = strings.ToLower(strings.ReplaceAll(question.Title, " ", "-"))
	question.Slug = reg.ReplaceAllString(question.Slug, "")
	db.Create(&question)
	err = json.NewEncoder(w).Encode(question)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get question
	params := mux.Vars(r)
	slug := params["slug"]

	// Checks if question exists
	if !XExists(slug, "question") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Question Not Found"})
		log.Handler("info", "Question not found", err)
		return
	}

	db.Where("slug = ?", slug).Find(&question)

	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	// Check if logged in user created the question
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	// Update the question
	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("warning", "JSON decoder error", err)
	db.Save(&question)

	// Return the question details
	err = json.NewEncoder(w).Encode(question)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	// Check if question exists
	params := mux.Vars(r)
	slug := params["slug"]

	if !XExists(slug, "question") {
		// Checks if question exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Question Not Found"})
		log.Handler("info", "Question not found", err)
		return
	}

	// Check if logged in user has permission to delete question
	db.Where("slug = ?", slug).Find(&question)
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}
	db.Where("slug = ?", slug).Delete(&question)
	w.WriteHeader(http.StatusNoContent)
	log.Handler("info", "Question deleted", nil)
}

func GetQuestionUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	db.Where("questionupvote_question_slug = ?", slug).Find(&questionUpVotes)
	err := json.NewEncoder(w).Encode(questionUpVotes)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostQuestionUpVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if XExists(slug, "qUpvote") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Upvote already Found"})
		log.Handler("info", "Upvote already found", err)
		return
	}

	db.Where("slug = ?", slug).First(&question)
	db.Find(&user, "email = ?", strings.ToLower(email))
	questionUpVote = QuestionUpVote{
		Question: question,
		User:     user,
	}
	db.Create(&questionUpVote)
	return
}

func DeleteQuestionUpvote(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	params := mux.Vars(r)
	slug, _ := params["slug"]

	db.Find(&user, "email = ?", strings.ToLower(email))
	db.Where("questionupvote_question_slug = ?", slug).Where(
		"questionupvote_user_id = ?", user.ID).Find(&questionUpVote)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != questionUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	db.Where("questionupvote_question_slug = ?", slug).Where(
		"questionupvote_user_id = ?", user.ID).Delete(&questionUpVote)
	w.WriteHeader(http.StatusNoContent)
	return
}

// Answers and Answer up votes

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if !XExists(slug, "answer") {
		// Checks if answer exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Answer Not Found"})
		log.Handler("info", "Answer not found", err)
		return
	}
	db.Where("slug = ?", slug).First(&answer)
	err := json.NewEncoder(w).Encode(answer)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func GetAnswers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&answers)
	err := json.NewEncoder(w).Encode(answers)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostAnswer(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&answer)
	log.Handler("warning", "JSON decoder error", err)
	db.Find(&user, "email = ?", strings.ToLower(email))
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	answer.Slug = strings.ToLower(strings.ReplaceAll(answer.Question.Title + "answer"+ strconv.Itoa(int(answer.ID)),
		" ", "-"))
	answer.Slug = reg.ReplaceAllString(answer.Slug, "")
	answer.User = user
	db.Create(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]

	// Checks if answer exists
	if !XExists(slug, "answer") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Answer Not Found"})
		log.Handler("info", "Answer not found", err)
		return
	}

	// Get answer
	db.Where("slug = ?", slug).Find(&answer)

	// Check if logged in user has permission to update answer
	if email != answer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&answer)
	log.Handler("warning", "JSON decoder error", err)
	db.Find(&user, "email = ?", strings.ToLower(email))
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	answer.Slug = strings.ToLower(strings.ReplaceAll(answer.Question.Title + "answer"+ strconv.Itoa(int(answer.ID)),
		" ", "-"))
	answer.Slug = reg.ReplaceAllString(answer.Slug, "")
	answer.User = user
	db.Save(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]

	// Checks if answer exists
	if !XExists(slug, "answer") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Answer Not Found"})
		log.Handler("info", "Answer not found", err)
		return
	}
	// Get answer
	db.Where("slug = ?", slug).Find(&answer)

	// Check if logged in user has permission to update answer
	if email != answer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	db.Where("slug = ?", slug).Delete(&answer)
	w.WriteHeader(http.StatusNoContent)
	log.Handler("info", "Answer deleted", nil)
	return
}

func GetAnswerUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	db.Where("answerupvote_answer_slug = ?", slug).Find(&answerUpVotes)
	err := json.NewEncoder(w).Encode(answerUpVotes)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostAnswerUpVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	params := mux.Vars(r)
	slug, _ := params["slug"]

	// Check if upvote exists
	if XExists(slug, "aUpvote") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Upvote already exists"})
		log.Handler("info", "Upvote already exists", err)
		return
	}

	db.Where("slug = ?", slug).First(&answer)
	db.Find(&user, "email = ?", strings.ToLower(email))
	answerUpVote = AnswerUpvote{
		Answer: answer,
		User:   user,
	}
	db.Create(&answerUpVote)
	return
}

func DeleteAnswerUpvote(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(w, r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	// Get question up vote
	params := mux.Vars(r)
	idInUint, _ := strconv.ParseUint(params["qid"], 10, 64)
	answerID := uint(idInUint)
	db.Find(&user, "email = ?", strings.ToLower(email))
	db.Where("questionupvote_question_id = ?", answerID).Where(
		"questionupvote_user_id = ?", user.ID).Find(&questionUpVote)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != questionUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Unauthorized"})
		log.Handler("warning", "JSON encoder error", err)
		return
	}

	db.Find(&user, "email = ?", strings.ToLower(email))
	db.Where("answerupvote_answer_id = ?", answerID).Where(
		"answerupvote_user_id = ?", user.ID).Delete(&answerUpVote)
	w.WriteHeader(http.StatusNoContent)
	return
}
