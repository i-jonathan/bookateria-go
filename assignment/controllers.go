package assignment

import (
	"bookateria-api-go/account"
	"bookateria-api-go/core"
	"bookateria-api-go/log"
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

type QuestionRequest struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Deadline        string `json:"deadline"`
	QuestionSlug    string `json:"question_slug"`
	SubmissionCount int    `json:"submission_count"`
}

func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, email := core.GetTokenEmail(r)

	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
		return
	}

	var questionR QuestionRequest
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := json.NewDecoder(r.Body).Decode(&questionR)
	log.Handler(err)

	deadline, _ := time.Parse(time.RFC3339, questionR.Deadline)

	// Remove all symbols and spaces to generate slug.
	regex, _ := regexp.Compile("[^a-zA-Z0-9]+")
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
	log.Handler(err)
	return
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if !XExists(slug, "problem") {
		// Checks if assignment problem exists
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	db.Preload(clause.Associations).First(&problem, slug)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(problem)
	log.Handler(err)
	return
}

func GetQuestions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch all assignment problems
	db.Preload(clause.Associations).Find(&problems)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(problems)
	log.Handler(err)
	return
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the logged in user
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
		return
	}

	// Get problem slug from url
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if problem exists
	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)

	// Check if user has permission to edit. Meaning, did the logged in use create this?
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&problem)
	log.Handler(err)
	db.Save(&problem)

	err = json.NewEncoder(w).Encode(problem)
	log.Handler(err)
	return
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
		return
	}

	// Get slug
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if problem exists
	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	//db.Find(&problem, "slug = ?", slug)
	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)
	// Check if logged in user is the creator
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
	}

	db.Where("slug = ?", slug).Delete(&problem)
	w.WriteHeader(http.StatusNoContent)
	return
}

// Endpoints for Submissions
func PostSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]

	if !XExists(questionSlug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}
	db.Preload(clause.Associations).Where("slug = ?", questionSlug).Find(&problem)
	_, email := core.GetTokenEmail(r)
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Handler(err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		return
	}

	var count int64

	db.Preload(clause.Associations).Model(&Submission{}).Where("user_email = ?", email).Count(&count)

	if int(count) == problem.SubmissionCount {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourHundred)
		log.Handler(err)
		return
	}
	fileNameExtension := strings.Split(header.Filename, ".")

	filename := fileNameExtension[0] + "_" + strconv.Itoa(int(count+1)) + fileNameExtension[1]

	sess := core.ConnectAWS()
	status, slug, err := core.S3Upload(sess, file, filename)
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		log.Handler(err)
		err = json.NewEncoder(w).Encode(core.FiveHundred)
		log.Handler(err)
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
	log.Handler(err)
	return
}

func GetSubmissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	slug := params["qSlug"]

	_, email := core.GetTokenEmail(r)
	fmt.Println(slug, email)

	if !XExists(slug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", slug).Find(&problem)
	fmt.Println(problem.ID, problem.User.Email)
	//db.Preload(clause.Associations).Find(&problem, "where slug = ?", slug)
	if email != problem.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.Handler(err)
		return
	}

	db.Preload(clause.Associations).Where("problem_id = ?", problem.ID).Find(&submissions)
	//db.Find(&submissions, "where problem_id = ?", problem.ID)
	err := json.NewEncoder(w).Encode(submissions)
	log.Handler(err)
	return
}

func GetSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]
	submissionSlug := params["sSlug"]

	if !XExists(questionSlug, "problem") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	if !XExists(submissionSlug, "submission") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.Handler(err)
		return
	}

	db.Preload(clause.Associations).Where("slug = ?", submissionSlug).Find(&submissions)

	err := json.NewEncoder(w).Encode(submission)
	log.Handler(err)
	return
}
