package documents

import (
	"bookateria-api-go/core"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var documents []Document


func GetDocuments(w http.ResponseWriter, _ *http.Request) {
	documents = append(documents,
		Document{
			ID:      "1",
			Title:   "Jay",
			Author:  "Jonathan",
			Summary: "I rock",
		})
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(documents)
	core.ErrHandler(err)
	return
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	documents = append(documents,
		Document{
			ID:      "1",
			Title:   "Jay",
			Author:  "Jonathan",
			Summary: "I rock",
		})
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range documents {
		if item.ID == params["id"] {
			err := json.NewEncoder(w).Encode(item)
			core.ErrHandler(err)
			return
		}
	}
	return
}

func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newDocument Document
	err := json.NewDecoder(r.Body).Decode(&newDocument)
	core.ErrHandler(err)
	newDocument.ID = strconv.Itoa(len(documents) + 1)

	documents = append(documents, newDocument)
	err = json.NewEncoder(w).Encode(newDocument)
	core.ErrHandler(err)
	return
}

func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, document := range documents {
		if document.ID == params["id"] {
			documents = append(documents[:i], documents[i+1:]...)
			var newDocument Document
			err := json.NewDecoder(r.Body).Decode(&newDocument)
			core.ErrHandler(err)
			newDocument.ID = params["id"]
			documents = append(documents, newDocument)
			err = json.NewEncoder(w).Encode(newDocument)
			core.ErrHandler(err)
			return
		}
	}
}

func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, document := range documents {
		if document.ID == params["id"] {
			documents = append(documents[:i], documents[i+1:]...)
			return
		}
	}
	err := json.NewEncoder(w).Encode(documents)
	core.ErrHandler(err)
}
