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
	answerUpVotes   []answerUpvote
	oneAUpVote      answerUpvote
	answers         []answer
	oneAnswer       answer
	db              = InitDatabase()
	oneQuestion     question
	questions       []question
	oneQUpVote      questionUpVote
	questionUpVotes []questionUpVote
	//questionTags    []questionTag
	user            account.User
)

// GetQuestion responds with a oneQuestion if the given slug exists
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slugToFind:= params["slug"]

	if !XExists(slugToFind, "question") {
		// Checks if oneQuestion exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).First(&oneQuestion, slugToFind)
	err := json.NewEncoder(w).Encode(oneQuestion)
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
	err := json.NewDecoder(r.Body).Decode(&oneQuestion)
	log.ErrorHandler(err)

	oneQuestion.Title = strings.Join(strings.Fields(oneQuestion.Title), " ")

	isValid := validator([]string{oneQuestion.Title})

	for _, tag := range oneQuestion.QuestionTags {
		if tag.Name == "" {
			isValid = false
		}
	}
	
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}
	
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
	oneQuestion.User = user

	// generate slug from title
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	oneQuestion.Slug = strings.Join(strings.Fields(oneQuestion.Title), " ")
	oneQuestion.Slug = strings.ToLower(strings.ReplaceAll(oneQuestion.Slug, " ", "-"))
	oneQuestion.Slug = reg.ReplaceAllString(oneQuestion.Slug, "")

	oneQuestion.Title = strings.Title(oneQuestion.Title)
	
	db.Create(&oneQuestion)
	err = json.NewEncoder(w).Encode(oneQuestion)
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
	
	// Get oneQuestion
	params := mux.Vars(r)
	slug := params["slug"]

	// Checks if oneQuestion exists
	if !XExists(slug, "question") {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Where("slug = ?", slug).Find(&oneQuestion)

	// Check if logged in user created the oneQuestion
	if email != oneQuestion.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// Update the oneQuestion
	err := json.NewDecoder(r.Body).Decode(&oneQuestion)
	log.ErrorHandler(err)

	isValid := validator([]string{oneQuestion.Title})

	for _, tag := range oneQuestion.QuestionTags {
		if tag.Name == "" {
			isValid = false
		}
	}
	
	if !isValid {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}

	db.Save(&oneQuestion)

	// Return the oneQuestion details
	err = json.NewEncoder(w).Encode(oneQuestion)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// DeleteQuestion removes an already created oneQuestion
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

	// Check if oneQuestion exists
	params := mux.Vars(r)
	slug := params["slug"]

	if !XExists(slug, "question") {
		// Checks if oneQuestion exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	// Check if logged in user has permission to delete oneQuestion
	db.Where("slug = ?", slug).Find(&oneQuestion)
	if email != oneQuestion.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	db.Where("slug = ?", slug).Delete(&oneQuestion)
	w.WriteHeader(http.StatusNoContent)
}

// GetQuestionUpVotes gets all the oneQuestion up votes
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

// PostQuestionUpVote creates a new upvote for a particular oneQuestion
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

	db.Where("slug = ?", slug).First(&oneQuestion)
	db.Find(&user, "email = ?", strings.ToLower(email))
	oneQUpVote = questionUpVote{
		Question: oneQuestion,
		User:     user,
	}
	db.Create(&oneQUpVote)
	log.AccessHandler(r, 200)
	return
}

// DeleteQuestionUpvote removes an upvote from a oneQuestion
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
		"questionupvote_user_id = ?", user.ID).Find(&oneQUpVote)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != oneQUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("questionupvote_question_slug = ?", slug).Where(
		"questionupvote_user_id = ?", user.ID).Delete(&oneQUpVote)
	w.WriteHeader(http.StatusNoContent)
	return
}

// Answers and oneAnswer up votes

// GetAnswer responds with an oneAnswer bt the slug given
func GetAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]

	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// Checks if oneAnswer exists
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	db.Find(&oneQuestion, "slug = ?", questionSlug)
	db.Where("slug = ?", slug).Where("question_id = ?", oneQuestion.ID).First(&oneAnswer)
	err := json.NewEncoder(w).Encode(oneAnswer)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetAnswers responds with all answers on a desired oneQuestion
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
	db.Find(&oneQuestion, "slug = ?", questionSlug)
	
	db.Find(&answers, "question_id = ?", oneQuestion.ID)
	err := json.NewEncoder(w).Encode(answers)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostAnswer for creating a new oneAnswer
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
	
	// Check if oneQuestion exists
	if !XExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	
	db.Find(&oneQuestion, "slug = ?", questionSlug)

	err := json.NewDecoder(r.Body).Decode(&oneAnswer)
	log.ErrorHandler(err)

	oneAnswer.Question = oneQuestion

	db.Find(&user, "email = ?", strings.ToLower(email))
	oneAnswer.User = user

	// Generate slug
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	slugText := strings.Join(strings.Fields(oneQuestion.Title), " ")
	oneAnswer.Slug = strings.ToLower(strings.ReplaceAll(slugText+"oneAnswer"+strconv.Itoa(int(oneAnswer.ID)), " ", "-"))
	oneAnswer.Slug = reg.ReplaceAllString(oneAnswer.Slug, "")

	db.Create(&oneAnswer)
	err = json.NewEncoder(w).Encode(oneAnswer)
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

	// Checks if oneAnswer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&oneQuestion, "slug = ?", questionSlug)
	// Get oneAnswer
	db.Where("slug = ?", slug).Where("question_id = ?", oneQuestion.ID).Find(&oneAnswer)

	// Check if logged in user has permission to update oneAnswer
	if email != oneAnswer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&oneAnswer)
	log.ErrorHandler(err)
	db.Find(&user, "email = ?", strings.ToLower(email))
	// reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	// oneAnswer.Slug = strings.ToLower(strings.ReplaceAll(oneAnswer.oneQuestion.Title+"oneAnswer"+strconv.Itoa(int(oneAnswer.ID)),
	// 	" ", "-"))
	oneAnswer.Slug = slug
	oneAnswer.User = user
	db.Save(&oneAnswer)
	err = json.NewEncoder(w).Encode(oneAnswer)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// DeleteAnswer removes a created oneAnswer
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

	// Checks if oneQuestion and oneAnswer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	// Get oneAnswer
	db.Find(&oneQuestion, "slug = ?", questionSlug)
	db.Where("slug = ?", slug).Where("question_id = ?", oneQuestion.ID).Find(&oneAnswer)

	// Check if logged in user has permission to update oneAnswer
	if email != oneAnswer.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("slug = ?", slug).Where("question_id = ?", oneQuestion.ID).Delete(&oneAnswer)
	w.WriteHeader(http.StatusNoContent)
	log.AccessHandler(r, 204)
	return
}

// GetAnswerUpVotes returns all upvotes on a given oneQuestion
func GetAnswerUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug := params["slug"]
	questionSlug := params["questionSlug"]
	
	// Checks if oneQuestion and oneAnswer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&oneQuestion, "slug = ?", questionSlug)
	db.Find(&oneAnswer, "question_id = ? AND slug = ?", oneQuestion.ID, slug)

	db.Preload(clause.Associations).Where("answer_id = ?", oneAnswer.ID).Find(&answerUpVotes)
	err := json.NewEncoder(w).Encode(answerUpVotes)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// PostAnswerUpVote upvotes an oneAnswer
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

	// Checks if oneQuestion and oneAnswer exists
	if !(XExists(questionSlug, "question") && XExists(slug, "answer")) {
		// If it doesn't return message accordingly
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Find(&oneQuestion, "slug = ?", questionSlug)
	db.Find(&oneAnswer, "question_id = ? AND slug = ?", oneQuestion.ID, slug)
	db.Find(&user, "email = ?", strings.ToLower(email))
	
	// Check if upvote exists
	var count int64
	db.Model(&answerUpvote{}).Where("user_id = ?", user.ID).Count(&count)

	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		log.AccessHandler(r, 409)
		return
	}

	oneAUpVote = answerUpvote{
		Answer: oneAnswer,
		User:   user,
	}
	db.Create(&oneAUpVote)
	log.AccessHandler(r, 200)
	return
}

// DeleteAnswerUpvote removes upvote from oneAnswer
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

	// Get oneQuestion up vote
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

	db.Find(&oneQuestion, "slug = ?", questionSlug)
	db.Find(&oneAnswer, "slug = ? AND question_id = ?", slug, oneQuestion.ID)

	db.Find(&user, "email = ?", strings.ToLower(email))
	
	// Check if upvote exists
	var count int64
	db.Model(&answerUpvote{}).Where("user_id = ? AND answer_id = ?", user.ID, oneAnswer.ID).Count(&count)

	if count > 0 {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		log.AccessHandler(r, 409)
		return
	}

	db.Find(&oneAUpVote, "user_id = ? AND answer_id = ?", user.ID, oneAnswer.ID)

	// Check if logged in user posted the upvote. If not, no permission to delete.
	if email != oneQUpVote.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Find(&oneAUpVote, "user_id = ? AND answer_id = ?", user.ID, oneAnswer.ID).Delete(&oneAUpVote)
	w.WriteHeader(http.StatusNoContent)
	log.AccessHandler(r, 204)
	return
}

// Search and queries
//  oneQuestion search

// QuestionSearch : Search for oneQuestion with query parameter
func QuestionSearch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Common oneQuestion words
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
	
	var questionList []question
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

// FilterQuestionByTags : Get oneQuestion that have a particular tag or tags
func FilterQuestionByTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filterQuery := r.URL.Query().Get("filter")
	tags := strings.Split(filterQuery, ",")
	var questionIDs []uint

	var questionTags []questionTag
	db.Where("name IN ?", tags).Find(&questionTags)

	for _, questionTag := range questionTags {
		questionIDs = append(questionIDs, questionTag.QuestionID)
	}
	
	db.Preload(clause.Associations).Find(&questions, "id IN ?", questionIDs)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(questions)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}