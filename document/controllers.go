package document

import (
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"strconv"
)

var (
	documents []Document
	document Document
	db = InitDatabase()
)

type Response struct {
	Message string `json:"message"`
}

func GetDocuments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Load data from DB
	db.Find(&documents) 
	err := json.NewEncoder(w).Encode(documents)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	documentID := params["id"]
	documentExists, _ := FilterBy("id", documentID)

	// Check If The Document Exists
	if !documentExists {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The document doesn't exist"})
		log.Handler("info", "The Document To Be Updated Does Not Exists", err)
		return

	}

	db.First(&document, documentID)
	err := json.NewEncoder(w).Encode(document)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.Handler("warning", "JSON decoder error", err)

	// Check If The Document Is A Duplicate Before Processing It
	isDuplicate := CheckDuplicate(&document)
	if isDuplicate {
		// The Document Is A Duplicate
		// Duplicate Documents Not Allowed

		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(Response{Message: "The document is a duplicate"})
		log.Handler("info", "Duplicate Document Detected", err)
		return
	}

	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.Handler("warning", "JSON decoder error", err)
	documentID := fmt.Sprint(document.ID)
	documentExists, _ := FilterBy("id", documentID)

	// Check If The Document Exists
	if !documentExists {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The document doesn't exist"})
		log.Handler("info", "The Document To Be Updated Does Not Exists", err)
		return

	}

	db.Save(&document)
	err = json.NewEncoder(w).Encode(document)
	log.Handler("warning", "JSON encoder error", err)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)
	db.Where("document_id = ?", idToDelete).Delete(&Document{})
	w.WriteHeader(http.StatusNoContent)
	log.Handler("info", "Document deleted", nil)
}