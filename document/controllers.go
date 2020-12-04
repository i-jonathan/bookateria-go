package document

import (
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var (
	documents []Document
	document  Document
	db        = InitDatabase()
)

type Response struct {
	Message string `json:"message"`
}

func GetDocuments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Load data from DB
	db.Find(&documents) 
	err := json.NewEncoder(w).Encode(documents)
	log.ErrorHandler(err)
	return
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	documentID := params["id"]
	ID, _ := strconv.ParseUint(documentID, 10, 0)

	// Check If The Document Exists
	if !DocumentExists(uint(ID)) {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The document doesn't exist"})
		log.ErrorHandler(err)
		return

	}

	db.First(&document, documentID)
	err := json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
	return
}

func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.ErrorHandler(err)

	// Check If The Document Is A Duplicate Before Processing It
	isDuplicate := CheckDuplicate(&document)
	if isDuplicate {
		// The Document Is A Duplicate
		// Duplicate Documents Not Allowed

		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(Response{Message: "The document is a duplicate"})
		log.ErrorHandler(err)
		return
	}

	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
	return
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.ErrorHandler(err)
	params := mux.Vars(r)
	idToUpdate, _ := strconv.ParseUint(params["id"], 10, 0)
	//documentID := strconv.FormatUint(uint64(document.ID), 10)

	// Check If The Document Exists
	if !DocumentExists(uint(idToUpdate)) {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The document doesn't exist"})
		log.ErrorHandler(err)
		return

	}

	// Checks If A Document Like That Already Exists
	if CheckDuplicate(&document) {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(Response{Message: "The document is a duplicate"})
		log.ErrorHandler(err)
		return
	}

	db.Save(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)

	//Check If The Document to Delete Exists
	if !DocumentExists(idToDelete) {
		//Deletion Of Non-Existent Documents Is Not Permitted
		//Throw An Error

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(Response{Message: "The document doesn't exist"})
		log.ErrorHandler(err)
		return

	}
	
	db.Where("id = ?", idToDelete).Delete(&document)
	w.WriteHeader(http.StatusNoContent)
}
