package account

import (
	"bookateriago/core"
	"bookateriago/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	users       []User
	user        User
	db          = InitDatabase()
	viperConfig = core.ReadViper()
	redisDB, _  = strconv.Atoi(fmt.Sprintf("%s", viperConfig.Get("redis.database")))
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s", viperConfig.Get("redis.address")),
		//Password: fmt.Sprintf("%s", viperConfig.Get("redis.password")),
		DB: redisDB,
	})
	ctx = context.Background()
)

type OTP struct {
	Email string `json:"email"`
	Pin   string `json:"pin"`
}

type OTPRequest struct {
	Email string `json:"email"`
}

// AllUsers gets and returns a list of all users in the DB
func AllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&users)
	err := json.NewEncoder(w).Encode(users)
	log.ErrorHandler(err)
	log.AccessHandler(r.URL.Path + " - [200]")
	return
}

// GetUser returns a user by id. TODO change to by slug
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]
	db.First(&user, userID)
	err := json.NewEncoder(w).Encode(user)
	log.ErrorHandler(err)
	log.AccessHandler(r.URL.Path + " - [200]")
	return
}

// PostUser for creating a new user. Does all the checks.
func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&user)
	log.ErrorHandler(err)
	var (
		email         = user.Email
		lastName      = user.LastName
		userName      = user.UserName
		password      = user.Password
		firstName     = user.FirstName
		safeNames     bool
		safeEmail     = EmailValidator(email)
		safePassword  = PasswordValidator(password)
		similarToUser = SimilarToUser(firstName, lastName, userName, password)
	)

	firstName, lastName, email, safeNames = UserDetails(firstName, lastName, email)

	if !safeNames {
		// Some or all of the details in the body are empty
		//	All fields are required
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [422]")
		return
	}

	if !safeEmail {
		// Issue with Email
		//Email couldn't be verified  or invalid email
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [422]")
		return
	}

	if similarToUser {
		// Issue with Password
		// Password is similar to user information
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [422]")
		return
	}

	if !safePassword {
		// Issue with Password
		//	Password doesn't go through the validator successfully
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [422]")
		return
	}

	passwordHash, _ := GeneratePasswordHash(password)

	user = User{
		UserName:        userName,
		FirstName:       firstName,
		LastName:        lastName,
		Email:           email,
		IsAdmin:         false,
		Password:        passwordHash,
		LastLogin:       time.Time{},
		IsActive:        false,
		IsEmailVerified: false,
	}

	db.Create(&user)
	err = json.NewEncoder(w).Encode(user)
	log.ErrorHandler(err)

	// Create OTP to verify email by
	// OTP expires in 30 minutes
	// Stored in Redis with key new_user_otp_email
	verifiableToken := GenerateOTP()
	err = redisClient.Set(ctx, "new_user_otp_"+email, verifiableToken, 30*time.Minute).Err()
	log.ErrorHandler(err)

	payload := struct {
		Token string
	}{
		Token: verifiableToken,
	}

	var status bool

	status, err = core.SendEmailNoAttachment(email, "OTP for Verification", payload, "token.txt")
	if !status {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(core.FiveHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [500]")
		return
	}
	log.ErrorHandler(err)
	return
}

// VerifyEmail is used to verify emails and make sure they exists.
// Supplementary to the regex check and the MX lookup
func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data OTP
	err := json.NewDecoder(r.Body).Decode(&data)
	log.ErrorHandler(err)

	// Gets the user and checks if the mail is already verified
	db.Find(&user, "email = ?", strings.ToLower(data.Email))
	if user.IsEmailVerified {
		w.WriteHeader(http.StatusTeapot)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [418]")
		return
	}

	// Gets the OTP stored in redis
	var storedOTP string
	key := "new_user_otp_" + data.Email
	storedOTP, err = redisClient.Get(ctx, key).Result()
	log.ErrorHandler(err)

	// If the OTP is empty, or the key doesn't exist or the pin provided is incorrect,
	// the pin has either elapsed the 30 minutes given or just plain wrong
	// So they need to request a new one
	if storedOTP == "" || storedOTP != data.Pin{
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [401]")
		return
	}

	w.WriteHeader(http.StatusOK)
	user.IsEmailVerified = true
	db.Save(&user)
	log.AccessHandler(r.URL.Path + " - [200]")
	return

}

// RequestOTP : In case the OTP sent expires, users can request for a new OTP
func RequestOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		data      OTPRequest
		storedOTP string
	)

	err := json.NewDecoder(r.Body).Decode(&data)

	db.Find(&user, "email = ?", data.Email)
	key := "new_user_otp_" + data.Email
	storedOTP, err = redisClient.Get(ctx, key).Result()

	if storedOTP == "" {
		verifiableToken := GenerateOTP()
		err = redisClient.Set(ctx, key, verifiableToken, 30*time.Minute).Err()
		storedOTP = verifiableToken
	}

	payload := struct {
		Token string
	}{
		Token: storedOTP,
	}
	var status bool
	status, err = core.SendEmailNoAttachment(data.Email, "OTP for Verification", payload, "token.txt")
	if !status {
		log.ErrorHandler(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(core.FiveHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r.URL.Path + " - [500]")
		return
	}
	w.WriteHeader(http.StatusOK)
	log.AccessHandler(r.URL.Path + " - [200]")
	return

}
