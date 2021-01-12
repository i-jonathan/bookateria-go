package assignment

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	oneSubmission submission
	submissions   []submission
	oneProblem    problem
	problems      []problem
	user          account.User
	db            = initDatabase()
)

// questionRequest - Accepted structure for taking data from request body
type questionRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Deadline        string `json:"deadline"`
	SubmissionCount int    `json:"submission_count"`
}

// postQuestion for creating a new assignment question
func postQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, email := core.GetTokenEmail(r)

	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	var questionR questionRequest
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := json.NewDecoder(r.Body).Decode(&questionR)
	log.ErrorHandler(err)

	deadline, _ := time.Parse(time.RFC3339, questionR.Deadline)

	/// get random string
	slugCode, err := rand.Int(rand.Reader, big.NewInt(999))
	if err != nil {
		fmt.Println(err)
	}

	// Remove all symbols and spaces to generate slug.
	regex, err := regexp.Compile("[^a-zA-Z0-9-]+")
	log.ErrorHandler(err)
	questionR.Title = strings.Join(strings.Fields(questionR.Title), "-")
	processed := regex.ReplaceAllString(questionR.Title, "")
	log.ErrorHandler(err)

	oneProblem = problem{
		Title:           questionR.Title,
		Description:     questionR.Description,
		Deadline:        deadline,
		User:            user,
		Slug:            processed + "-" + slugCode.String(),
		SubmissionCount: questionR.SubmissionCount,
	}

	db.Create(&oneProblem)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(oneProblem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// getQuestion gets a question by the slug passed in, in the url
func getQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if !xExists(slug, "question") {
		// Checks if assignment problem exists
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).Find(&oneProblem, "slug = ?", slug)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(oneProblem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// getQuestions gets all assignment questions in the db.
func getQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch all assignment problems
	db.Preload(clause.Associations).Find(&problems)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(problems)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// updateQuestion adjusts an already existing assignment question
func updateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the logged in user
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// Get problem slug from url
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if problem exists
	if !xExists(slug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&oneProblem)

	// Check if user has permission to edit. Meaning, did the logged in use create this?
	if email != oneProblem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&oneProblem)
	log.ErrorHandler(err)
	db.Save(&oneProblem)

	err = json.NewEncoder(w).Encode(oneProblem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// deleteQuestion removes a question
func deleteQuestion(w http.ResponseWriter, r *http.Request) {
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

	// Get slug
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if problem exists
	if !xExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&oneProblem)
	// Check if logged in user is the creator
	if email != oneProblem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("slug = ?", slug).Delete(&oneProblem)
	w.WriteHeader(http.StatusNoContent)
	log.AccessHandler(r, 204)
	return
}

// Endpoints for Submissions

// PostSubmission creates a new submission for the question which is identified by the slug
func PostSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]

	if !xExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	db.Preload(clause.Associations).Where("slug = ?", questionSlug).Find(&oneProblem)
	_, email := core.GetTokenEmail(r)
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.ErrorHandler(err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}

	var count int64

	db.Preload(clause.Associations).Model(&submission{}).Where("user_id = ? and question_id = ?", user.ID, oneProblem.ID).Count(&count)
	fmt.Println(count)
	fmt.Println(oneProblem.SubmissionCount)
	if int(count) >= oneProblem.SubmissionCount {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}
	fileNameExtension := strings.Split(header.Filename, ".")

	filename := fileNameExtension[0] + "_" + strconv.Itoa(int(count+1)) + fileNameExtension[1]

	status, slug, err := core.S3Upload(file, filename)
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		log.ErrorHandler(err)
		err = json.NewEncoder(w).Encode(core.FiveHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 500)
		return
	}

	submissionSlug := oneProblem.Slug + "-" + strings.Join(strings.Fields(oneProblem.User.Alias), "-") + "-" + strconv.Itoa(int(count) + 1)

	oneSubmission = submission{
		Problem:     oneProblem,
		User:        user,
		FileSlug:    slug,
		Slug:        submissionSlug,
		Submissions: count + 1,
	}

	db.Create(&oneSubmission)
	err = json.NewEncoder(w).Encode(oneSubmission)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// getSubmissions returns all submissions to the individual who created the question
func getSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	slug := params["qSlug"]

	_, email := core.GetTokenEmail(r)

	if !xExists(slug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&oneProblem)
	//db.Preload(clause.Associations).Find(&problem, "where slug = ?", slug)
	if email != oneProblem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Preload(clause.Associations).Where("problem_id = ?", oneProblem.ID).Find(&submissions)
	//db.Find(&submissions, "where problem_id = ?", problem.ID)
	err := json.NewEncoder(w).Encode(submissions)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// getSubmission returns the submission of a particular person
func getSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]
	submissionSlug := params["aSlug"]

	_, email := core.GetTokenEmail(r)

	if !xExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", submissionSlug).Find(&oneSubmission)

	if email != oneSubmission.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewEncoder(w).Encode(oneSubmission)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}
