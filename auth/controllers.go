package auth

import (
	"bookateria-api-go/account"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var (
	jwtKey = []byte("Hello")
	user account.User
	db = account.InitDatabase()
	cred Credentials
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

type Claim struct {
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

	claims := Claim{
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

	_ = json.NewEncoder(w).Encode(TokenResponse{
		Name:   "Token",
		Value:  tokenString,
		Expiry: expirationTime,
	})
	return
}