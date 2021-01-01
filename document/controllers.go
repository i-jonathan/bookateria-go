package document

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm/clause"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	db = InitDatabase()
)

//GetDocuments fetches all documents in the database
//
func GetDocuments(w http.ResponseWriter, _ *http.Request) {
	var documents []Document
	w.Header().Set("Content-Type", "application/json")
	// Load data from DB
	db.Preload(clause.Associations).Find(&documents)
	err := json.NewEncoder(w).Encode(documents)
	log.ErrorHandler(err)
}

//GetDocument fetches a specific document from the database
func GetDocument(w http.ResponseWriter, r *http.Request) {
	var document Document
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
	var (
		document Document
		tags     []Tag
		tag      Tag
		email    string
		user     account.User
	)

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

	//Check for user attached to mail
	db.Find(&user, "email = ?", strings.ToLower(email))
	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")

	//If The Regexp Doesn't Compile, Throw An Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	//Validate The Title Field
	title, err := validate(r.FormValue("title"))

	//If The Title Field Is Not Valid Throw An Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	//Validate The Author Field
	author, err := validate(r.FormValue("author"))

	//If The Author Field Is Not Valid, Throw An Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	//Default Value For Edition
	edition := 0

	fmt.Println(r.FormValue("edition"))

	//Check If The Edition Sent By The User Is Not Empty
	if r.FormValue("edition") != "" {
		var err error
		edition, err = strconv.Atoi(r.FormValue("edition"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(core.FourONine)
			log.ErrorHandler(err)
			return
		}
	}

	//Get tags from request and split into a slice
	strTags := strings.Split(r.FormValue("tags"), ",")

	//Parse tags and store into the Tags model
	for _, strTag := range strTags {
		tag.TagName = strings.TrimSpace(string(strTag))
		tag.Slug = strings.ReplaceAll(strings.ToLower(tag.TagName), " ", "-")
		tag.Slug = reg.ReplaceAllString(tag.Slug, "")
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
	slug := strings.ToLower(strings.ReplaceAll(document.Title+"-"+document.Author+"-"+fmt.Sprint(edition), " ", "-"))
	slug = reg.ReplaceAllString(slug, "")
	document.Slug = slug

	//Create an entry for the document in the database
	db.Create(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

//UpdateDocument overwrites the details of a specified document with the provided ones.
func UpdateDocument(w http.ResponseWriter, r *http.Request) {
	var (
		document Document
		temp     Document
		email    string
		user     account.User
	)

	w.Header().Set("Content-Type", "application/json")

	//Checks If Current User Is Logged In
	if _, email = core.GetTokenEmail(r); email == "" {
		fmt.Println("login")
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		return
	}

	//Decode The Request Into A Temporary Document Variable
	err := json.NewDecoder(r.Body).Decode(&temp)
	log.ErrorHandler(err)
	params := mux.Vars(r)

	//Parse The ID To Be Updated
	idToUpdate, err := strconv.ParseUint(params["id"], 10, 0)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	// Check If The Document Exists
	if !xExists(uint(idToUpdate)) {
		// If The Document Doesn't Exist
		// Users Shouldn't Be Allowed To Modify What Doesn't Exists

		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(core.FourOFour)
		log.ErrorHandler(err)
		return

	}

	//Gets The Document With The Specified ID
	db.Find(&document, "id = ?", idToUpdate)

	//Gets The Uploader Of The Document Gotten Above
	db.Find(&user, "id = ?", document.UploaderID)

	//Check If The Person Updating Is Authorized To Do So.
	if email != user.Email {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		return
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9-]+")

	//If The Regexp Doesn't Compile, Throw An Error
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		return
	}

	//Check If The Title Is To Be Updated
	if string(temp.Title) != "" {
		title, err := validate(string(temp.Title))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(core.FourTwoTwo)
			log.ErrorHandler(err)
			return
		}
		document.Title = title
	}

	//Check If The Author Is To Be Updated
	if string(temp.Author) != "" {

		//Validate The Title Field
		author, err := validate(string(temp.Author))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(core.FourTwoTwo)
			log.ErrorHandler(err)
			return
		}
		document.Author = author
	}

	//Check If The Edition Is To Be Updated
	if document.Edition != 0 {

		//Check If Edition Is An Integer
		edition, err := strconv.Atoi(fmt.Sprint(temp.Edition))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(core.FourTwoTwo)
			log.ErrorHandler(err)
			return
		}

		document.Edition = edition
	}

	//Check If Tags Are Empty
	if len(temp.Tags) > 0 {

		//Parse Tags If They're Not Empty
		for _, tag := range temp.Tags {
			slug := strings.ReplaceAll(strings.ToLower(string(tag.TagName)), " ", "-")
			tag.Slug = reg.ReplaceAllString(slug, "")
			document.Tags = append(document.Tags, tag)
		}
	}

	//Gets document slug from request and stores it into the document
	slug := strings.ToLower(strings.ReplaceAll(document.Title+"-"+document.Author+"-"+fmt.Sprint(document.Edition), " ", "-"))
	slug = reg.ReplaceAllString(slug, "")
	document.Slug = slug

	//Update The Document Uploader
	document.Uploader = user

	//Save The Document
	db.Save(&document)
	err = json.NewEncoder(w).Encode(document)
	log.ErrorHandler(err)
}

//DeleteDocument removes a specified document from the DB
func DeleteDocument(w http.ResponseWriter, r *http.Request) {
	var (
		document Document
		tags     []Tag
	)

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

	//Load The Tags Associated With The Document To Delete
	db.Find(&tags, "document_id = ?", idToDelete)

	for _, tag := range tags {
		db.Delete(&tag)
	}

	db.Where("id = ?", idToDelete).Delete(&document)
	w.WriteHeader(http.StatusNoContent)
}
