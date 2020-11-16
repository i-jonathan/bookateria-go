package core

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"net/http"
	"net/smtp"
)

type Claims struct {
	Email string `json:"email"`
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
		//w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	token, _ := jwt.ParseWithClaims(authorization, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := token.Claims.(*Claims)
	if ok && token.Valid {
	} else {
		//w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	email := claims.Email

	return token, email
}

func SendEmail(to, from, subject, messageBody string) bool {
	viperConfig := ReadViper()
	fromMail := fmt.Sprintf("%s", viperConfig.Get("email.address"))
	password := fmt.Sprintf("%s", viperConfig.Get("email.password"))

	toMail := []string{
		to,
	}

	smtpHost := fmt.Sprintf("%s", viperConfig.Get("email.host"))
	smtpPort := fmt.Sprintf("%d", viperConfig.Get("email.port"))
	auth := smtp.PlainAuth("", fromMail, password, smtpHost)

	//message := fmt.Sprintf("From: %s\r\n To: %s\r\n Subject: %s\r\n\r\n %s\r\n", from, to, subject, messageBody)
	message := "From: "+ from +"\r\n " +
		"To:" + to + "\r\n" +
		"Subject:" + subject + "\r\n" +
		"\r\n" +
		messageBody + ".\r\n"

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, toMail, []byte(message))
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Sent")
	return true
}