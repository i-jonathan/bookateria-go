package auth

import (
	"bookateriago/account"
	"bookateriago/core"
	"bookateriago/log"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	viperConfig = core.ReadViper()
	jwtKey      = []byte(fmt.Sprintf("%s", viperConfig.Get("settings.key")))
	db          = account.InitDatabase()
	redisDb, _  = strconv.Atoi(fmt.Sprintf("%d", viperConfig.Get("redis.database")))
	ctx         = context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", viperConfig.Get("redis.address")),
		// Password: fmt.Sprintf("%s", viperConfig.Get("redis.password")),
		DB:       redisDb,
	})
	user account.User
	cred credentials
)

// tokenResponse is the structure of the access token
type tokenResponse struct {
	Name   string    `json:"name"`
	Value  string    `json:"value"`
	Expiry time.Time `json:"expiry"`
}

// credentials is the expected struct for the sign in endpoint
type credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
	StayIn   bool   `json:"stay_in"`
}

// tokenClaims for building jwt token
type tokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// SignIn takes a post request with the credentials to be logged in with
func signIn(w http.ResponseWriter, r *http.Request) {
	// Reads the body for email and password, gets the user and the password from DB
	// Compares the password, if correct, returns the token
	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.ErrorHandler(err)
		log.AccessHandler(r, 400)
		return
	}
	db.Find(&user, "email = ?", strings.ToLower(cred.Email))
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}
	expectedPassword := user.Password
	correct, _ := account.ComparePassword(cred.Password, expectedPassword)

	if !correct {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	redisTime := 168 * time.Hour

	expirationTime := time.Now().Add(168 * time.Hour)

	if cred.StayIn {
		expirationTime = time.Now().Add(720 * time.Hour)
		redisTime = 720 * time.Hour
	}

	claims := tokenClaims{
		Email: cred.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errToken := token.SignedString(jwtKey)
	if errToken != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.ErrorHandler(err)
		log.AccessHandler(r, 500)
		return
	}
	err = redisClient.Set(ctx, user.Email, tokenString, redisTime).Err()
	log.ErrorHandler(err)

	err = json.NewEncoder(w).Encode(tokenResponse{
		Name:   "Authorization",
		Value:  tokenString,
		Expiry: expirationTime,
	})

	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}

// func RefreshToken(w http.ResponseWriter, r *http.Request) {
// 	// This function refreshes the token of a sign in function. Once the expiration time is within
// 	// 30 seconds, it send back a new token
// 	// Should also save the new token to redis
// 	authorization := r.Header.Get("Authorization")
// 	if authorization == "" {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	claims := &tokenClaims{}

// 	token, err := jwt.ParseWithClaims(authorization, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	if !token.Valid {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	claims.ExpiresAt = expirationTime.Unix()
// 	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	_ = json.NewEncoder(w).Encode(tokenResponse{
// 		Name:   "Token",
// 		Value:  tokenString,
// 		Expiry: expirationTime,
// 	})
// 	return
// }

// Logout immediately deletes a valid
func logout(w http.ResponseWriter, r *http.Request) {
	// This function deletes the token from redis, rendering it invalid
	// Make sure to delete the token on the frontend too

	w.Header().Set("Content-Type", "application/json")
	_, email := core.GetTokenEmail(r)
	if email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		err := json.NewEncoder(w).Encode(core.FourOOne)
		log.ErrorHandler(err)
		log.AccessHandler(r, 401)
		return
	}

	redisClient.Del(ctx, email)
	err := json.NewEncoder(w).Encode(core.TwoHundred)
	log.ErrorHandler(err)
	log.AccessHandler(r, 200)
	return
}
