package forum

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
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

// GetQuestion responds with a question if the given slug exists
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slugToFind:= params["slug"]

	if !XExists(slugToFind, "question") {
		// Checks if question exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).First(&question, slugToFind)
	err := json.NewEncoder(w).Encode(question)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetQuestions gets a list of all questions in the database
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Preload(clause.Associations).Find(&questions)
	err := json.NewEncoder(w).Encode(questions)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostQuestion is the function that handles creation of a new questions
func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&question)
	log.ErrorHandler(err)
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	
	// get user
	db.Find(&user, "email = ?", strings.ToLower(email))
	question.User = user

	// generate slug from title
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	question.Slug = strings.Join(strings.Fields(question.Title), " ")
	question.Slug = strings.ToLower(strings.ReplaceAll(question.Slug, " ", "-"))
	question.Slug = reg.ReplaceAllString(question.Slug, "")
	
	db.Create(&question)
	err = json.NewEncoder(w).Encode(question)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// UpdateQuestion adjusts data of already created questions
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	
	// Get question
	params := mux.Vars(r)
	slug := params["slug"]

	// Checks if question exists
	if !XExists(slug, "question") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Where("slug = ?", slug).Find(&question)

	// Check if logged in user created the question
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// Update the question
	err := json.NewDecoder(r.Body).Decode(&question)
	log.ErrorHandler(err)
	db.Save(&question)

	// Return the question details
	err = json.NewEncoder(w).Encode(question)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// DeleteQuestion removes an already created question
func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// Check if question exists
	params := mux.Vars(r)
	slug := params["slug"]

	if !XExists(slug, "question") {
		// Checks if question exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	// Check if logged in user has permission to delete question
	db.Where("slug = ?", slug).Find(&question)
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	db.Where("slug = ?", slug).Delete(&question)
	w.WriteHeader(http.StatusNoContent)
}

// GetQuestionUpVotes get's all the question upvotes
func GetQuestionUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]

	db.Where("questionupvote_question_slug = ?", slug).Find(&questionUpVotes)
	err := json.NewEncoder(w).Encode(questionUpVotes)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostQuestionUpVote creates a new upvote for a particular question
func PostQuestionUpVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	params := mux.Vars(r)
	slug := params["slug"]

	if XExists(slug, "qUpvote") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Where("slug = ?", slug).First(&question)
	db.Find(&user, "email = ?", strings.ToLower(email))
	questionUpVote = QuestionUpVote{
		Question: question,
		User:     user,
	}
	db.Create(&questionUpVote)
	log.AccessHandler(r, 200)
	return
}

// DeleteQuestionUpvote removes an upvote from a question
func DeleteQuestionUpvote(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]

	db.Find(&user, "email = ?", strings.ToLower(email))
	db.Where("questionupvote_question_slug = ?", slug).Where(
		"questionupvote_user_id = ?", user.ID).Find(&questionUpVote)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != questionUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("questionupvote_question_slug = ?", slug).Where(
		"questionupvote_user_id = ?", user.ID).Delete(&questionUpVote)
	w.WriteHeader(http.StatusNoContent)
	return
}

// Answers and Answer up votes

// GetAnswer responds with an answer bt the slug given  
func GetAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// Checks if answer exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	db.Find(&question, "slug = ?", questionSlug)
	db.Where("slug = ?", slug).Where("question_id = ?", question.ID).First(&answer)
	err := json.NewEncoder(w).Encode(answer)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetAnswers responds with all answers on a desired question
func GetAnswers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	questionSlug := params["questionSlug"]

	if !XExists(questionSlug, "question") {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	db.Find(&question, "slug = ?", questionSlug)
	
	db.Find(&answers, "question_id = ?", question.ID)
	err := json.NewEncoder(w).Encode(answers)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostAnswer for creating a new answer
func PostAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	params := mux.Vars(r)
	questionSlug := params["questionSlug"]
	
	// Check if question exists
	if !XExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	
	db.Find(&question, "slug = ?", questionSlug)

	err := json.NewDecoder(r.Body).Decode(&answer)
	log.ErrorHandler(err)

	answer.Question = question

	db.Find(&user, "email = ?", strings.ToLower(email))
	answer.User = user

	// Generate slug
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	slugText := strings.Join(strings.Fields(question.Title), " ")
	answer.Slug = strings.ToLower(strings.ReplaceAll(slugText+"answer"+strconv.Itoa(int(answer.ID)), " ", "-"))
	answer.Slug = reg.ReplaceAllString(answer.Slug, "")

	db.Create(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// UpdateAnswer endpoint for updating answers
func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	// Checks if answer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&question, "slug = ?", questionSlug)
	// Get answer
	db.Where("slug = ?", slug).Where("question_id = ?", question.ID).Find(&answer)

	// Check if logged in user has permission to update answer
	if email != answer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&answer)
	log.ErrorHandler(err)
	db.Find(&user, "email = ?", strings.ToLower(email))
	// reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	// answer.Slug = strings.ToLower(strings.ReplaceAll(answer.Question.Title+"answer"+strconv.Itoa(int(answer.ID)),
	// 	" ", "-"))
	answer.Slug = slug
	answer.User = user
	db.Save(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// DeleteAnswer removes a created answer
func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	// Checks if question and answer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	// Get answer
	db.Find(&question, "slug = ?", questionSlug)
	db.Where("slug = ?", slug).Where("question_id = ?", question.ID).Find(&answer)

	// Check if logged in user has permission to update answer
	if email != answer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("slug = ?", slug).Where("question_id = ?", question.ID).Delete(&answer)
	w.WriteHeader(http.StatusNoContent)
	log.AccessHandler(r, 204)
	return
}

// GetAnswerUpVotes returns all upvotes on a given question
func GetAnswerUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]
	
	// Checks if question and answer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&question, "slug = ?", questionSlug)
	db.Find(&answer, "question_id = ? AND slug = ?", question.ID, slug)

	db.Preload(clause.Associations).Where("answer_id = ?", answer.ID).Find(&answerUpVotes)
	err := json.NewEncoder(w).Encode(answerUpVotes)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostAnswerUpVote upvotes an answer
func PostAnswerUpVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	// Checks if question and answer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&question, "slug = ?", questionSlug)
	db.Find(&answer, "question_id = ? AND slug = ?", question.ID, slug)
	db.Find(&user, "email = ?", strings.ToLower(email))
	
	// Check if upvote exists
	var count int64
	db.Model(&AnswerUpvote{}).Where("user_id = ?", user.ID).Count(&count)

	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		log.AccessHandler(r, 409)
		return
	}

	answerUpVote = AnswerUpvote{
		Answer: answer,
		User:   user,
	}
	db.Create(&answerUpVote)
	log.AccessHandler(r, 200)
	return
}

// DeleteAnswerUpvote removes upvote from answer
func DeleteAnswerUpvote(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// Get question up vote
	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&question, "slug = ?", questionSlug)
	db.Find(&answer, "slug = ? AND question_id = ?", slug, question.ID)

	db.Find(&user, "email = ?", strings.ToLower(email))
	
	// Check if upvote exists
	var count int64
	db.Model(&AnswerUpvote{}).Where("user_id = ? AND answer_id = ?", user.ID, answer.ID).Count(&count)

	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		log.AccessHandler(r, 409)
		return
	}

	db.Find(&answerUpVote, "user_id = ? AND answer_id = ?", user.ID, answer.ID)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != questionUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Find(&answerUpVote, "user_id = ? AND answer_id = ?", user.ID, answer.ID).Delete(&answerUpVote)
	w.WriteHeader(http.StatusNoContent)
	log.AccessHandler(r, 204)
	return
}

// Search and queries
//  Question search

// QuestionSearch : Search for question withd query parameter
func QuestionSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Common question words
	questionWords := []string{"why", "who", "what", "how", "whom", "when", "where", "are", "is", "the", "whose"}

	searchTerm := r.URL.Query().Get("search")
	regex, _ := regexp.Compile("[\\W+]")
	final := regex.ReplaceAllString(searchTerm, "")

	for _, word := range questionWords {
		final = strings.ReplaceAll(strings.ToLower(final), word, "")
	}

	if final == "" {
		final = regex.ReplaceAllString(searchTerm, "")
	}

	individualWords := strings.Fields(final)
	
	var questionList []Question
	for _, word := range individualWords {
		db.Where("lower(title) LIKE ?", "%"+strings.ToLower(word)+"%").Find(&questions)
		questionList = append(questionList, questions...)
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(questionList)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// FilterQuestionByTags : Get question that have a particular tag or tags  
func FilterQuestionByTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filterQuery := r.URL.Query().Get("filter")
	tags := strings.Split(filterQuery, ",")
	var questionIDs []uint

	var questionTags []QuestionTag
	db.Where("name IN ?", tags).Find(&questionTags)

	for _, questionTag := range questionTags {
		questionIDs = append(questionIDs, questionTag.QuestionID)
	}
	
	db.Preload(clause.Associations).Find(&questions, "id IN ?", questionIDs)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(questions)
	log.AccessHandler(r, 200)
	return
}