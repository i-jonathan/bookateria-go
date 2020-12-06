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

type Response struct {
	Message string `json:"message"`
}

func AllUsers(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db.Find(&users)
	err := json.NewEncoder(w).Encode(users)
	log.ErrorHandler(err)
	log.ErrorHandler(err)
	return
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	userID := params["id"]
	db.First(&user, userID)
	err := json.NewEncoder(w).Encode(user)
	log.ErrorHandler(err)

	return
}

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
		err := json.NewEncoder(w).Encode(Response{Message: "Name and Email are required"})
		log.ErrorHandler(err)
		return
	}

	if !safeEmail {
		// Issue with Email
		//Email couldn't be verified  or invalid email
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Incorrect Email"})
		log.ErrorHandler(err)
		return
	}

	if similarToUser {
		// Issue with Password
		// Password is similar to user information
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Password is similar to user info"})
		log.ErrorHandler(err)
		return
	}

	if !safePassword {
		// Issue with Password
		//	Password doesn't go through the validator successfully
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(Response{Message: "Unsafe Password"})
		log.ErrorHandler(err)
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

	// Create OTP to verify email by
	// OTP expires in 30 minutes
	// Stored in Redis with key new_user_otp_email
	verifiableToken := GenerateOTP()
	err = redisClient.Set(ctx, "new_user_otp_"+email, verifiableToken, 30*time.Minute).Err()
	if err != nil {
		//	Do stuff
	}

	payload := struct {
		Token string
	}{
		Token: verifiableToken,
	}

	var status bool

	status, err = core.SendEmailNoAttachment(email, "OTP for Verification", payload, "token.txt")
	if !status {
		// TODO Log error.
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Response{Message: "Email not sent. Server Error"})
		log.ErrorHandler(err)
	}
	log.ErrorHandler(err)
	return
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data OTP
	err := json.NewDecoder(r.Body).Decode(&data)
	log.ErrorHandler(err)

	// Gets the user and checks if the mail is already verified
	db.Find(&user, "email = ?", strings.ToLower(data.Email))
	if user.IsEmailVerified {
		w.WriteHeader(http.StatusTeapot)
		_ = json.NewEncoder(w).Encode(Response{
			Message: "Email already Verified",
		})
	}

	// Gets the OTP stored in redis
	var storedOTP string
	key := "new_user_otp_" + data.Email
	storedOTP, err = redisClient.Get(ctx, key).Result()
	log.ErrorHandler(err)

	// If the OTP is empty, or the key doesn't exist, the pin has most likely elapsed the 30 minutes given
	// So they need to request a new one
	if storedOTP == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(Response{Message: "Your OTP might have expired. Request a new one"})
		log.ErrorHandler(err)
		return
	}

	// If the pin exists, and it's not the same as that provided by the client, unauthorized error is raised
	if storedOTP != data.Pin {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(Response{Message: "Pin is invalid. Please check your mail and try again"})
		return
	}

	w.WriteHeader(http.StatusOK)
	user.IsEmailVerified = true
	db.Save(&user)
	err = json.NewEncoder(w).Encode(Response{Message: "Your Email has been Verified. You can now Login"})
	return

}

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
		// TODO Log error.
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Response{Message: "Email not sent. Server Error"})
		log.ErrorHandler(err)
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Response{Message: "OTP has been sent to mail"})
	log.ErrorHandler(err)
	return

}
