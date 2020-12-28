package document

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"net/http"
	//"regexp"
	"strconv"
	"strings"
)

var (
	documents []Document
	tags      []Tag
	tag       Tag
	document  Document
	db        = InitDatabase()
	user      account.User
	email     string
)

//GetDocuments fetches all documents in the database
//
func GetDocuments(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Load data from DB
	db.Preload(clause.Associations).Find(&documents)
	err := json.NewEncoder(w).Encode(documents)
	log.ErrorHandler(err)
}

//GetDocument fetches a specific document from the database
func GetDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	documentID := params["id"]
	ID, _ := strconv.ParseUint(documentID, 10, 0)

	// Check If The Document Exists
	if !xExists(uint(ID)) {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		return

	}

	db.Preload(clause.Associations).Find(&document, "id = ?", documentID)
	err := json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

//PostDocument puts a provided document into the db
func PostDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")

	//Checks If Current User Is Logged In
	if _, email = core.GetTokenEmail(r); email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		return
	}

	//Creates memory space to store form-data
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.ErrorHandler(err)
	}

	db.Find(&user, "email = ?", strings.ToLower(email)) //Check for user attached to mail
	//reg, _ := regexp.Compile("[^a-zA-Z0-9]+")

	title, author, edition, err := validate(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	strTags := strings.Split(r.FormValue("tags"), ",") //Get tags from request and split into a slice

	//Parse tags and store into the Tags model
	for _, strTag := range strTags {
		tag.TagName = strings.TrimSpace(string(strTag))
		tag.Slug = strings.ReplaceAll(strings.ToLower(string(strTag)), " ", "-")
		tags = append(tags, tag)
	}

	//Store documents info
	document = Document{
		Title:    title,
		Author:   author,
		Edition:  edition,
		Tags:     tags,
		Summary:  r.FormValue("summary"),
		Uploader: user,
	}

	//Checks if the document is a duplicate
	if checkDuplicate(&document) {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		return
	}

	//Gets document slug from request and stores it into the document
	slug := strings.ToLower(strings.ReplaceAll(document.Title+"-"+document.Author+"-"+r.FormValue("edition"), " ", "-"))
	document.Slug = slug

	//Create an entry for the document in the database
	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

//UpdateDocument overwrites the details of a specified document with the provided ones.
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&document)
	log.ErrorHandler(err)
	params := mux.Vars(r)
	idToUpdate, _ := strconv.ParseUint(params["id"], 10, 0)
	//documentID := strconv.FormatUint(uint64(document.ID), 10)

	// Check If The Document Exists
	if !xExists(uint(idToUpdate)) {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		return

	}

	// Checks If A Document Like That Already Exists
	if checkDuplicate(&document) {
		w.WriteHeader(http.StatusConflict)
		err := json.NewEncoder(w).Encode(core.FourONine)
		log.ErrorHandler(err)
		return
	}

	db.Save(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

//DeleteDocument removes a specified document from the DB
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInUint, _ := strconv.ParseUint(id, 10, 64)
	idToDelete := uint(idInUint)

	//Check If The Document to Delete Exists
	if !xExists(idToDelete) {
		//Deletion Of Non-Existent Documents Is Not Permitted
		//Throw An Error

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		return

	}

	db.Where("id = ?", idToDelete).Delete(&document)
	w.WriteHeader(http.StatusNoContent)
}
