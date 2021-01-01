package core

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	template2 "html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var viperConfig = ReadViper()
var key = fmt.Sprintf("%s", viperConfig.Get("email.key"))

// parseTemplate is for preparing the template from the email/templates directory
func parseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("email/templates/%s", templateFileName))

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
func SendEmailNoAttachment(toMail, subject string, data interface{}, template string) (bool, error) {
	emailBody, err := parseTemplate(template, data)

	if err != nil {
		return false, err
	}

	from := mail.NewEmail("Bookateria", "noreply@bookateria.net")
	to := mail.NewEmail("Me", toMail)
	message := mail.NewSingleEmail(from, subject, to, emailBody, emailBody)
	// a := mail.NewAttachment()
	// a.SetContent("TG9yZW0gaXBzdW0gZG9sb3Igc2l0IGFtZXQsIGNvbnNlY3RldHVyIGFkaXBpc2NpbmcgZWxpdC4gQ3JhcyBwdW12")
	// a.SetFilename("config.yaml")
	// a.SetType("text/plain")
	// a.SetDisposition("attachment")
	
	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)

	if response.StatusCode != 202 {
		return false, err
	}
	
	return true, nil
}


// SendEmailWithAttachment for attaching files to an email.
func SendEmailWithAttachment(toMail, subject, content, fileDir, fileName, template string, data interface{}) (bool, error) {
	fileBytes, err := ioutil.ReadFile(fileDir + fileName)

	if err != nil {
		return false, err
	}

	emailBody, err := parseTemplate(template, data)

	if err != nil {
		return false, err
	}

	fileType := http.DetectContentType(fileBytes)
	fileData := base64.StdEncoding.EncodeToString(fileBytes)

	from := mail.NewEmail("Bookateria", "hi@bookateria.net")
	to := mail.NewEmail("Me", toMail)
	// htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, emailBody, emailBody)
	a := mail.NewAttachment()
	a.SetContent(fileData)
	a.SetFilename(fileDir+fileName)
	a.SetType(fileType)
	a.SetDisposition("attachment")
	
	message.AddAttachment(a)
	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)
	if response.StatusCode != 202 {
		return false, err
	}
	return true, nil
}
