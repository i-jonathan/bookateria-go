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
)

var (
	submission Submission
	submissions []Submission
	question  Question
	questions []Question
	user      account.User
	db        = InitDatabase()
)

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

	db.Find(&user, "email = ?", strings.ToLower(email))

	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("info", "Couldn't decode Body", err)
	question.User = user

	// Remove all symbols and spaces to generate slug.
	regex, _ := regexp.Compile("[^a-zA-Z0-9]+")
	processed := regex.ReplaceAllString(question.Title, "")
	slug := strings.ReplaceAll(processed, " ", "-")
	//randInt, _ := rand.Int(rand.Reader, big.NewInt(9999))
	question.QuestionSlug = slug

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

func PostSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	params := mux.Vars(r)
	questionSlug := params["qSlug"]

	if !XExists(questionSlug, "question") {
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message:"Question not found"})
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