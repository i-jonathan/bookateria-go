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

func S3Upload(sess *session.Session, file multipart.File, filename string) (bool, error) {
	// This takes a file from a multipart/form and uploads to an AWS S3 bucket
	// Pass in the session from the ConnectAWS function, the file from the multipart form
	// and the filename from header.filename

	bucketName := fmt.Sprintf("%s", viperConfig.Get("aws.bucket"))

	uploader := s3manager.NewUploader(sess)
	_, err := uploader.Upload(&s3manager.UploadInput{
		ACL: 	aws.String("public-read"),
		Body:   file,
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}