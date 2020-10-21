package auth

import (
	"bookateria-api-go/account"
	"bookateria-api-go/core"
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
	ctx 		= context.Background()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", viperConfig.Get("redis.address")),
		Password: fmt.Sprintf("%s", viperConfig.Get("redis.password")),
		DB:       redisDb,
	})
	user        account.User
	cred        Credentials
)

type TokenResponse struct {
	Name	string 		`json:"name"`
	Value	string		`json:"value"`
	Expiry	time.Time 	`json:"expiry"`
}

type Response struct {
	Message	string `json:"message"`
}

type Credentials struct {
	Password	string	`json:"password"`
	Email		string	`json:"email"`
}

type Claims struct {
	Email	string	`json:"email"`
	jwt.StandardClaims
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// Reads the body for email and password, gets the user and the password from DB
	// Compares the password, if correct, returns the token
	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	db.Find(&user, "email = ?", strings.ToLower(cred.Email))
	if user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(Response{Message: "User Not found"})
		return
	}
	expectedPassword := user.Password
	correct, _ := account.ComparePassword(cred.Password, expectedPassword)

	if !correct {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(Response{Message: "Incorrect Details"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := Claims{
		Email:          cred.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, errToken := token.SignedString(jwtKey)
	if errToken != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = redisClient.Set(ctx, user.Email, tokenString, 5*time.Minute).Err()
	if err != nil {
		panic(err)
	}

	_ = json.NewEncoder(w).Encode(TokenResponse{
		Name:   "Token",
		Value:  tokenString,
		Expiry: expirationTime,
	})
	redisT, err := redisClient.Get(ctx, user.Email).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		panic(err)
	}

	fmt.Println(redisT)
	return
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	// This function refreshes the token of a sign in function. Once the expiration time is within
	// 30 seconds, it send back a new token
	// Should also save the new token to redis
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(authorization, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(TokenResponse{
		Name:   "Token",
		Value:  tokenString,
		Expiry: expirationTime,
	})
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// This function deletes the token from redis, rendering it invalid
	// Make sure to delete the token on the frontend too

	w.Header().Set("Content-Type", "application/json")
	_, email := GetTokenEmail(w, r)
	redisClient.Del(ctx, email)
	_ = json.NewEncoder(w).Encode(Response{Message: "Successfully logged out"})
	return
}

func GetTokenEmail(w http.ResponseWriter, r *http.Request) (*jwt.Token, string) {
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	token, _ := jwt.ParseWithClaims(authorization, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	email := claims.Email

	return token, email
}