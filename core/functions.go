package core

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"net/http"
)

type Claims struct {
	Email	string	`json:"email"`
	jwt.StandardClaims
}

func ReadViper() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.SetConfigFile("config.yaml")
	_ = viperConfig.ReadInConfig()

	return viperConfig
}

func GetTokenEmail(w http.ResponseWriter, r *http.Request) (*jwt.Token, string) {
	authorization := r.Header.Get("Authorization")
	viperConfig := ReadViper()
	jwtKey := []byte(fmt.Sprintf("%s", viperConfig.Get("settings.key")))
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