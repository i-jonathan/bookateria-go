package core

import (
	"bookateriago/log"
	"bytes"
	"errors"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	template2 "html/template"
	"path/filepath"
)

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

	client := sendgrid.NewSendClient(key)
	response, err := client.Send(message)

	log.ErrorHandler(err)

	if response.StatusCode != 202 {
		return false, err
	}

	return true, nil
}

// SendEmailWithAttachment for attaching files to an email.
/*
func SendEmailWithAttachment(toMail, subject, fileDir, fileName, template string, data interface{}) (bool, error) {
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
	log.ErrorHandler(err)
	if response.StatusCode != 202 {
		return false, err
	}
	return true, nil
}
*/
