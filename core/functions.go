package core

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"mime/multipart"
	"net/http"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// ReadViper: A simple function utilizing the viper package for reading from configuration file.
// Reads specifically from config.yaml located in the root directory.
func ReadViper() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.SetConfigFile("config.yaml")
	_ = viperConfig.ReadInConfig()

	return viperConfig
}

// GetTokenEmail is used to get the token as well as email address of logged in users
// It returns both the token and the email if the user is logged in and the token is valid
// Else it returns nil and an empty string: "".
// Reads the token from the request header and breaks it down to get the user.
func GetTokenEmail(r *http.Request) (*jwt.Token, string) {
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

func ConnectAWS() *session.Session {
	viperConfig := ReadViper()
	// This is used to connect to AWS with the credentials
	accessKeyID := fmt.Sprintf("%s", viperConfig.Get("aws.accessKeyID"))
	secretAccessKey := fmt.Sprintf("%s", viperConfig.Get("aws.secretAccessKey"))
	bucketRegion := fmt.Sprintf("%s", viperConfig.Get("aws.region"))

	sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
			Region:      aws.String(bucketRegion),
		})

	if err != nil {
		fmt.Println(err)
	}

	return sess
}

// S3Upload takes a file from a multipart/form and uploads to an AWS S3 bucket
// Pass in the session from the ConnectAWS function, the file from the multipart form
// and the filename from header.filename
// Returns true, the slug, and nil if successful
func S3Upload(sess *session.Session, file multipart.File, filename string) (bool, string, error) {

	bucketName := fmt.Sprintf("%s", viperConfig.Get("aws.bucket"))

	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		ACL:    aws.String("public-read"),
		Body:   file,
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		return false, "", err
	}

	fileSlug := "https://" + bucketName + "." + "s3.amazonaws.com/" + filename
	return true, fileSlug, nil
}
