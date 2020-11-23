package assignment

import (
	"bookateria-api-go/account"
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	submission  Submission
	submissions []Submission
	question    Question
	questions   []Question
	user        account.User
	db          = InitDatabase()
)

type QuestionRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
	QuestionSlug string `json:"question_slug"`
}

type Response struct {
	Message string
}

func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, email := core.GetTokenEmail(r)

	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "User not logged in."})
		log.Handler("info", "Json Encoder ish", err)
		return
	}

	var questionR QuestionRequest
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := json.NewDecoder(r.Body).Decode(&questionR)
	log.Handler("info", "Couldn't decode Body", err)

	deadline, _ := time.Parse(time.RFC3339, questionR.Deadline)

	// Remove all symbols and spaces to generate slug.
	regex, _ := regexp.Compile("[^a-zA-Z0-9]+")
	processed := regex.ReplaceAllString(question.Title, "")
	slug := strings.ReplaceAll(processed, " ", "-")
	//randInt, _ := rand.Int(rand.Reader, big.NewInt(9999))

	question = Question{
		Title:       questionR.Title,
		Description: questionR.Description,
		Deadline:    deadline,
		User:        user,
		Slug:        slug,
	}

	db.Create(&question)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(question)
	log.Handler("info", "JSON Encoder", err)
	return
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	slug, _ := params["slug"]

	if !XExists(slug, "question") {
		// Checks if assignment question exists
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Assignment question not found"})
		log.Handler("info", "JSON Encoder error", err)
		return
	}

	db.Preload(clause.Associations).First(&question, slug)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(question)
	log.Handler("info", "JSON Encoder again", err)
	return
}

func GetQuestions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Fetch all assignment questions
	db.Preload(clause.Associations).Find(&questions)
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(questions)
	log.Handler("info", "json encoder", err)
	return
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the logged in user
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("info", "JSON Encoder", err)
		return
	}

	// Get question slug from url
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if question exists
	if !XExists(slug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The requested resource couldn't be located in this timeline"})
		log.Handler("info", "JSON Encoder again", err)
		return
	}

	db.Find(&question, "slug = ?", slug)

	// Check if user has permission to edit. Meaning, did the logged in use create this?
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "You do not have permissions to access here"})
		log.Handler("info", "JSON Encoder", err)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("info", "JSON Decoder", err)
	db.Save(&question)

	err = json.NewEncoder(w).Encode(question)
	log.Handler("info", "Really tired of doing this", err)
	return
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if user is logged in
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "Login Required"})
		log.Handler("info", "JSON Encoder", err)
		return
	}

	// Get slug
	params := mux.Vars(r)
	slug := params["slug"]

	// Check if question exists
	if !XExists(slug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Resource not found"})
		log.Handler("info", "JSON Encoder", err)
		return
	}

	db.Find(&question, "slug = ?", slug)
	// Check if logged in user is the creator
	if email != question.User.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(Response{Message: "You do not have that permissions"})
		log.Handler("info", "JSON Encoder", err)
	}

	db.Where("slug = ?", slug).Delete(&question)
	w.WriteHeader(http.StatusNoContent)
	return
}

// Endpoints for Submissions
func PostSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]

	if !XExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "Question not found"})
		log.Handler("info", "Json Encoder", err)
		return
	}
	db.Find(&question, "slug = ?", questionSlug)
	_, email := core.GetTokenEmail(r)
	db.Find(&user, "email = ?", strings.ToLower(email))

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Handler("info", "Something about parsing multipart form", err)
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Response{Message: "something went wrong retrieving the file"})
		return
	}

	filename := header.Filename

	sess := core.ConnectAWS()
	status, slug, err := core.S3Upload(sess, file, filename)
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		log.Handler("info", "S3Upload error", err)
		err = json.NewEncoder(w).Encode(Response{Message: "Couldn't upload file successfully"})
		log.Handler("info", "JSON Encoder error", err)
	}

	submission = Submission{
		Question: question,
		User:     user,
		FileSlug: slug,
	}

	db.Create(&submission)
	err = json.NewEncoder(w).Encode(submission)
	log.Handler("info", "JSON Encoder", err)
	return
}
