package account

import (
	"bookateria-api-go/log"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

var (
	users []User
	user User
	db = InitDatabase()
)

type Response struct {
	Message string `json:"message"`
}

func AllUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&users)
	err := json.NewEncoder(w).Encode(users)
	log.Handler("warning", "JSON encoder error", err)
	log.Handler("info", "All Users Endpoint returned values", nil)
	return
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]
	db.First(&user, userID)
	err := json.NewEncoder(w).Encode(user)
	log.Handler("warning", "JSON encoder error", err)

	return
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	log.Handler("warning", "JSON decoder error", err)
	var (
		email     = user.Email
		lastName  = user.LastName
		userName  = user.UserName
		password  = user.Password
		firstName = user.FirstName
		safeNames bool
		safeEmail = EmailValidator(email)
		safePassword = PasswordValidator(password)
		similarToUser = SimilarToUser(firstName, lastName, userName, password)
	)

	firstName, lastName, email, safeNames = UserDetails(firstName, lastName, email)

	if !safeNames {
		// Some or all of the details in the body are empty
		//	All fields are required
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Name and Email are required"})
		log.Handler("info", "Text in Body not accepted", err)
		return
	}

	if !safeEmail {
		// Issue with Email
		//Email couldn't be verified  or invalid email
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Incorrect Email"})
		log.Handler("info", "Wrong mail", err)
		return
	}

	if similarToUser {
		// Issue with Password
		// Password is similar to user information
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Password is similar to user info"})
		log.Handler("info", "Bad Password", err)
		return
	}

	if !safePassword {
		// Issue with Password
		//	Password doesn't go through the validator successfully
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Unsafe Password"})
		log.Handler("info", "Password is cheap", err)
		return
	}

	passwordHash, _ := GeneratePasswordHash(password)

	user = User{
		UserName:  userName,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsAdmin:   false,
		Password:  passwordHash,
		LastLogin: time.Time{},
		IsActive:  false,
	}

	db.Create(&user)
	err = json.NewEncoder(w).Encode(user)
	log.Handler("warning", "JSON encoder error", err)
	return
}
