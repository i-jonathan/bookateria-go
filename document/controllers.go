package document

import (
	"bookateria-api-go/core"
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
	core.ErrHandler(err)
	return
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	documentID := params["id"]
	db.First(&document, documentID)
	err := json.NewEncoder(w).Encode(document)
	core.ErrHandler(err)
	return
}

func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	core.ErrHandler(err)
	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	core.ErrHandler(err)
	return
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	core.ErrHandler(err)
	db.Save(&document)
	err = json.NewEncoder(w).Encode(document)
	core.ErrHandler(err)
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)
	db.Where("document_id = ?", idToDelete).Delete(&Document{})
	//err := json.NewEncoder(w).Encode(documents)
	//core.ErrHandler(err)
	w.WriteHeader(http.StatusNoContent)
}
