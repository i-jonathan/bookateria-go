package assignment

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	submission  Submission
	submissions []Submission
	problem     Problem
	problems    []Problem
	user        account.User
	db          = InitDatabase()
)

// QuestionRequest - Accepted structure for taking data from request body
type QuestionRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Deadline        string `json:"deadline"`
	SubmissionCount int    `json:"submission_count"`
}

// PostQuestion for creating a new assignment question
func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, email := core.GetTokenEmail(r)

	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	var questionR QuestionRequest
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := json.NewDecoder(r.Body).Decode(&questionR)
	log.ErrorHandler(err)

	deadline, _ := time.Parse(time.RFC3339, questionR.Deadline)

	// Remove all symbols and spaces to generate slug.
	problem.Title = strings.Join(strings.Fields(problem.Title), " ")
	regex, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	processed := regex.ReplaceAllString(problem.Title, "")
	slug := strings.ReplaceAll(processed, " ", "-")
	fmt.Println(slug)

	problem = Problem{
		Title:           questionR.Title,
		Description:     questionR.Description,
		Deadline:        deadline,
		User:            user,
		Slug:            slug,
		SubmissionCount: questionR.SubmissionCount,
	}

	db.Create(&problem)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(problem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetQuestion gets a question by the slug passed in, in the url
func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if !XExists(slug, "problem") {
		// Checks if assignment problem exists
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).First(&problem, slug)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(problem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetQuestions gets all assignment questions in the db.
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch all assignment problems
	db.Preload(clause.Associations).Find(&problems)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(problems)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// UpdateQuestion adjusts an already existing assignment question
func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
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
	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)

	// Check if user has permission to edit. Meaning, did the logged in use create this?
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&problem)
	log.ErrorHandler(err)
	db.Save(&problem)

	err = json.NewEncoder(w).Encode(problem)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// DeleteQuestion removes a question
func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
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
	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)
	// Check if logged in user is the creator
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Where("slug = ?", slug).Delete(&problem)
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

	if !XExists(questionSlug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}
	db.Preload(clause.Associations).Where("slug = ?", questionSlug).Find(&problem)
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

	db.Preload(clause.Associations).Model(&Submission{}).Where("user_email = ?", email).Count(&count)

	if int(count) == problem.SubmissionCount {
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

	submissionSlug := problem.Slug + problem.User.FirstName + "-" + problem.User.LastName

	submission = Submission{
		Problem:     problem,
		User:        user,
		FileSlug:    slug,
		Slug:        submissionSlug,
		Submissions: count + 1,
	}

	db.Create(&submission)
	err = json.NewEncoder(w).Encode(submission)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetSubmissions returns all submissions to the individual who created the question
func GetSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	slug := params["qSlug"]

	_, email := core.GetTokenEmail(r)

	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)
	fmt.Println(problem.ID, problem.User.Email)
	//db.Preload(clause.Associations).Find(&problem, "where slug = ?", slug)
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	db.Preload(clause.Associations).Where("problem_id = ?", problem.ID).Find(&submissions)
	//db.Find(&submissions, "where problem_id = ?", problem.ID)
	err := json.NewEncoder(w).Encode(submissions)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// GetSubmission returns the submission of a particular person
func GetSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]
	submissionSlug := params["sSlug"]

	_, email := core.GetTokenEmail(r)

	if !XExists(questionSlug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		log.AccessHandler(r, 404)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", submissionSlug).Find(&submission)

	if email != submission.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	err := json.NewEncoder(w).Encode(submission)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}
