package core

import (
	"bookateriago/log"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

type tokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	redisDB, _  = strconv.Atoi(fmt.Sprintf("%s", os.Getenv("redis_database")))
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s", os.Getenv("redis_address")),
		Password: fmt.Sprintf("%s", os.Getenv("redis_password")),
		DB: redisDB,
	})
	ctx = context.Background()
)

// ReadViper : A simple function utilizing the viper package for reading from configuration file.
// Reads specifically from config.yaml located in the root directory.
//func ReadViper() *viper.Viper {
//	viperConfig := viper.New()
//	viperConfig.SetConfigFile(".env")
//	_ = viperConfig.ReadInConfig()
//
//	return viperConfig
//}

// GetTokenEmail is used to get the token as well as email address of logged in users
// It returns both the token and the email if the user is logged in and the token is valid
// Else it returns nil and an empty string: "".
// Reads the token from the request header and breaks it down to get the user.
func GetTokenEmail(r *http.Request) (*jwt.Token, string) {
	authorization := r.Header.Get("Authorization")
	jwtKey := []byte(fmt.Sprintf("%s", os.Getenv("settings_key")))
	if authorization == "" {
		//w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	token, _ := jwt.ParseWithClaims(authorization, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := token.Claims.(*tokenClaims)
	if ok && token.Valid {
	} else {
		//w.WriteHeader(http.StatusUnauthorized)
		return nil, ""
	}

	email := claims.Email

	storedOTP, err := redisClient.Get(ctx, email).Result()
	log.ErrorHandler(err)

	if storedOTP == "" || storedOTP != authorization {
		return nil, ""
	}

	return token, email
}

// connectAWS connects to AWS with correct credentials and creates a session
func connectAWS() *session.Session {
	// This is used to connect to AWS with the credentials
	accessKeyID := fmt.Sprintf("%s", os.Getenv("aws_accessKeyID"))
	secretAccessKey := fmt.Sprintf("%s", os.Getenv("aws_secretAccessKey"))
	bucketRegion := fmt.Sprintf("%s", os.Getenv("aws_region"))

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
func S3Upload(file multipart.File, filename string) (bool, string, error) {

	sess := connectAWS()
	bucketName := fmt.Sprintf("%s", os.Getenv("aws_bucket"))

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


// ResponseData checks if a previous and next page exists for a certain endpoint
// returns too boolean values. Previous and next
func ResponseData(count int, r *http.Request) (int, bool, bool) {
	prev, next := false, false
	var page, pageSize int
	if count != 0 {
		var err error
		if r.URL.Query().Get("page") == "" {
			page = 1
		} else {
			page, err = strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				log.ErrorHandler(err)
			}
		}
		if r.URL.Query().Get("page_size") == "" {
			pageSize = 10
		} else {
			pageSize, err = strconv.Atoi(r.URL.Query().Get("page_size"))
			if err != nil {
				log.ErrorHandler(err)
			}
		}

		if (page * pageSize) < count {
			next = true
		}

		if page > 1 {
			prev = true
		}

	}

	return page, prev, next
}
