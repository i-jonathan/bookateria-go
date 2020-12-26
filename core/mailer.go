package core

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	template2 "html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"
)

var gmailService *gmail.Service
var viperConfig = ReadViper()

type data struct {
	ReceiverName string
	SenderName   string
}

// OAuthGmailService Creates a connection to gmail service
func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     fmt.Sprintf("%s", viperConfig.Get("email.client.id")),
		ClientSecret: fmt.Sprintf("%s", viperConfig.Get("email.client.secret")),
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("%s", viperConfig.Get("email.redirect")),
	}

	token := oauth2.Token{
		AccessToken:  fmt.Sprintf("%s", viperConfig.Get("token.access")),
		RefreshToken: fmt.Sprintf("%s", viperConfig.Get("token.refresh")),
		TokenType:    fmt.Sprintf("%s", viperConfig.Get("token.type")),
		Expiry:       time.Now(),
	}

	tokenSource := config.TokenSource(context.Background(), &token)

	service, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))

	if err != nil {
		fmt.Println(err)
	}

	gmailService = service
	if gmailService != nil {
		fmt.Println("Email Service Initiated")
	}
}

// parseTemplate is for preparing the template from the email/templates directory
func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("../email/templates/%s", templateFileName))

	if err != nil {
		return "", errors.New("invalid template name")
	}

	template, err := template2.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	buff := new(bytes.Buffer)

	if err = template.Execute(buff, data); err != nil {
		return "", err
	}

	body := buff.String()
	return body, nil
}

// SendEmailNoAttachment is for sending emails with no attachments, like OTP, password reset
func SendEmailNoAttachment(to, title string, data interface{}, template string) (bool, error) {
	emailBody, err := parseTemplate(template, data)

	if err != nil {
		return false, errors.New("unable to parse email template")
	}

	var message gmail.Message

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + title + "\n"
	mime := "MIME-versionL 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	msg := []byte(emailTo + subject + mime + "\n" + emailBody)

	message.Raw = base64.URLEncoding.EncodeToString(msg)

	_, err = gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}

	return true, nil
}

func randStr(strSize int) string {
	dictionary := "0123456789abcdefghijklmnopqrstwxyzABCDEFGHIJKLMNOPQRWXYZ"

	var strBytes = make([]byte, strSize)
	_, _ = rand.Read(strBytes)
	for key, value := range strBytes {
		strBytes[key] = dictionary[value%byte(len(dictionary))]
	}

	return string(strBytes)
}

func chunkSplit(body string, limit int, end string) string {
	var charSlice []rune

	// Add characters to charSlice
	for _, char := range body {
		charSlice = append(charSlice, char)
	}

	result := ""

	for len(charSlice) >= 1 {
		// convert slice to string
		// But insert end at the limit
		result = result + string(charSlice[:limit]) + end

		// discard elements that were copied over to result
		charSlice = charSlice[limit:]

		// change the limit to cater for the last few words
		if len(charSlice) < limit {
			limit = len(charSlice)
		}

	}
	return result
}

// SendEmailWithAttachment for attaching files to an email.
func SendEmailWithAttachment(to, subject, content, fileDir, fileName string) (bool, error) {
	var message gmail.Message

	fileBytes, err := ioutil.ReadFile(fileDir + fileName)

	if err != nil {
		return false, err
	}

	fileMIMEType := http.DetectContentType(fileBytes)
	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	boundary := randStr(32)

	messageBody := []byte("Content-Type: multipart/mixed; boundary=" + boundary + " \n" +
		"MIME-Version: 1.0\n" +
		"to: " + to + "\n" +
		"subject: " + subject + "\n\n" +

		"--" + boundary + "\n\n" +
		content + "\n\n" +
		"--" + boundary + "\n" +

		"Content-Type: " + fileMIMEType + "; name=" + string('"') + fileName + string('"') + " \n" +
		"MIME-Version: 1.0\n" +
		"Content-Transfer-Encoding: base64\n" +
		"Content-Disposition: attachment; filename=" + string('"') + fileName + string('"') + " \n\n" +
		chunkSplit(fileData, 76, "\n") +
		"--" + boundary + "--")

	message.Raw = base64.URLEncoding.EncodeToString(messageBody)

	// Send the mail
	_, err = gmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}

	return true, nil
}
