package account

import (
	"bookateriago/core"
	"bookateriago/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var (
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

// otp is the structure of the OTP itself
type otp struct {
	Email string `json:"email"`
	Pin   string `json:"pin"`
}

// otpRequest carries parameters for requesting OTPs
type otpRequest struct {
	Email string `json:"email"`
}

// allUsers gets and returns a list of all users in the DB
func allUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	db.Find(&users)
	err := json.NewEncoder(w).Encode(users)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// getUser returns a user by id. TODO change to by slug
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	params := mux.Vars(r)
	userID := params["id"]
	db.Find(&user, "id = ?", userID)
	err := json.NewEncoder(w).Encode(user)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// postUser for creating a new user. Does all the checks.
func postUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	log.ErrorHandler(err)
	var (
		email         = strings.ToLower(user.Email)
		alias      = user.Alias
		userName      = user.UserName
		password      = user.Password
		fullName     = user.FullName
		safeNames     bool
		safeEmail     = emailValidator(email)
		safePassword  = passwordValidator(password)
		similarToUser = similarToUser(fullName, alias, userName, password)
	)

	safeNames = userDetails(fullName, alias, userName)

	if safeNames {
		// Some or all of the details in the body are empty
		//	All fields are required
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r, 422)
		return
	}

	if !safeEmail {
		// Issue with Email
		//Email couldn't be verified  or invalid email
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r, 422)
		return
	}

	if similarToUser {
		// Issue with Password
		// Password is similar to user information
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r, 422)
		return
	}

	if !safePassword {
		// Issue with Password
		//	Password doesn't go through the validator successfully
		w.WriteHeader(http.StatusUnprocessableEntity)
		err := json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r, 422)
		return
	}

	passwordHash, err := generatePasswordHash(password)
	log.ErrorHandler(err)

	user = User{
		UserName:        userName,
		FullName:        fullName,
		Alias:           alias,
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
	verifiableToken := generateOTP()
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
		log.AccessHandler(r, 500)
		return
	}
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// verifyEmail is used to verify emails and make sure they exists.
// Supplementary to the regex check and the MX lookup
func verifyEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		data otp
		user User
	)
	err := json.NewDecoder(r.Body).Decode(&data)
	log.ErrorHandler(err)

	// Gets the user and checks if the mail is already verified
	db.Find(&user, "email = ?", strings.ToLower(data.Email))
	if user.IsEmailVerified {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
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
	if storedOTP == "" || storedOTP != data.Pin {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	w.WriteHeader(http.StatusOK)
	user.IsEmailVerified = true
	db.Save(&user)
	log.AccessHandler(r, 200)

	redisClient.Del(ctx, key)
	return

}

// requestOTP : In case the OTP sent expires, users can request for a new OTP
func requestOTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var (
		data      otpRequest
		storedOTP string
		user User
	)

	err := json.NewDecoder(r.Body).Decode(&data)

	db.Find(&user, "email = ?", data.Email)
	key := "new_user_otp_" + data.Email
	storedOTP, err = redisClient.Get(ctx, key).Result()

	if storedOTP == "" {
		verifiableToken := generateOTP()
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
		log.AccessHandler(r, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.AccessHandler(r, 200)
	return

}

// resetPasswordRequest handles the request to reset a password. Sends a mail to the user, containing the OTP
func resetPasswordRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body otpRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	log.ErrorHandler(err)
	
	// Verify email
	emailStatus := emailValidator(body.Email)
	if !emailStatus {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}

	var count int64
	var user User
	db.Find(&user, "email = ?", body.Email).Count(&count)

	if count <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// generate token
	var data otp
	data = otp{
		Email: body.Email,
		Pin: generateOTP(),
	}

	// save token to redis
	err = redisClient.Set(ctx, "password_reset_"+data.Email, data.Pin, 30*time.Minute).Err()
	if err != nil {
		log.ErrorHandler(err)
		w.WriteHeader(http.StatusInternalServerError)
		err := json.NewEncoder(w).Encode(core.FiveHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 500)
		return
	}

	// Send token to email
	payload := struct {
		Token string
	}{
		Token: data.Pin,
	}

	status, err := core.SendEmailNoAttachment(data.Email, "Reset Password", payload, "password_reset.txt")
	if !status {
		log.ErrorHandler(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(core.FiveHundred)
		log.ErrorHandler(err)
		log.AccessHandler(r, 500)
		return
	}

	// respond okay
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(core.TwoHundred)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// resetPassword actually resets the users password. Based on data provided and generated
func resetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// take email, token and new password
	body := struct {
		Email	 string	`json:"email"`
		OTP		 string	`json:"otp"`
		Password string	`json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&body)
	log.ErrorHandler(err)
	
	// check email for existence
	var user User
	err = db.Find(&user, "email = ?", body.Email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusUnauthorized)
		log.ErrorHandler(err)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// check if token exists
	storedOtp, err := redisClient.Get(ctx, "password_reset_"+body.Email).Result()
	log.ErrorHandler(err)
	
	if storedOtp != body.OTP {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	// validate password
	safePassword := passwordValidator(body.Password)
	if !safePassword {
		w.WriteHeader(http.StatusUnprocessableEntity)
		err = json.NewEncoder(w).Encode(core.FourTwoTwo)
		log.ErrorHandler(err)
		log.AccessHandler(r, 422)
		return
	}

	// Generate password hash and save
	hashedPassword, err := generatePasswordHash(body.Password)
	log.ErrorHandler(err)

	user.Password = hashedPassword
	db.Save(&user)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(core.TwoHundred)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)

	// Delete from redis
	redisClient.Del(ctx, "password_reset_"+body.Email)

	return
}