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
	document Document
	db = InitDatabase()
)

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
	db.First(&document, documentID)
	err := json.NewEncoder(w).Encode(document)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.Handler("warning", "JSON decoder error", err)
	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	log.Handler("warning", "JSON encoder error", err)
	return
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)

	log.Handler("warning", "JSON decoder error", err)
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