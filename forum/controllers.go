package forum

import (
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	answerUpVotes []AnswerUpvote
	answers []Answer
	answer Answer
	db = InitDatabase()
	question Question
	questions []Question
	questionUpVotes []QuestionUpVote
)

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idInUint, _ := strconv.ParseUint(params["id"], 10, 64)
	questionID := uint(idInUint)
	db.First(&question, questionID)
	err := json.NewEncoder(w).Encode(question)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func GetQuestions(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&questions)
	err := json.NewEncoder(w).Encode(questions)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("warning", "JSON decoder error", err)
	db.Create(&question)
	err = json.NewEncoder(w).Encode(question)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&question)
	log.Handler("warning", "JSON decoder error", err)
	db.Save(&question)
	err = json.NewEncoder(w).Encode(question)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)
	db.Delete(&Question{}, idToDelete)
	w.WriteHeader(http.StatusNoContent)
	log.Handler("info", "Question deleted", nil)
}

func GetQuestionUpVotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	idInUint, _ := strconv.ParseUint(params["id"], 10, 64)
	questionID := uint(idInUint)
	db.Where("questionupvote_question_id = ?", questionID).Find(&questionUpVotes)
	err := json.NewEncoder(w).Encode(questionUpVotes)
	log.Handler("info", "JSON Encoder error", err)
	return
}

func PostQuestionUpVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

}

// Answers

func GetAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	answerID := params["id"]
	db.First(&answer, answerID)
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
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&answer)
	log.Handler("warning", "JSON decoder error", err)
	db.Create(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateAnswer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&answer)
	log.Handler("warning", "JSON decoder error", err)
	db.Save(&answer)
	err = json.NewEncoder(w).Encode(answer)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)
	db.Delete(&Answer{}, idToDelete)
	w.WriteHeader(http.StatusNoContent)
	log.Handler("info", "Question deleted", nil)
}
