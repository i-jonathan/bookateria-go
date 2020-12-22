package log

import (
	"log"
	"net/http"
	"os"
	"strings"
)

func ErrorHandler(err error) {
	file, issue := os.OpenFile("log/error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if issue != nil {
		log.Printf("Error opening file: %v", issue)
		return
	}

	if err == nil {
		return
	}

	log.SetOutput(file)
	log.Println(err)
	err = file.Close()
	if err == nil {
		return
	}
	log.Println(err)
	return
}

func AccessHandler(r *http.Request, code int) {
	file, issue := os.OpenFile("log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if issue != nil {
		log.Printf("Error opening file: %v", issue)
		return
	}

	log.SetOutput(file)
	text := r.Method + " - " + r.URL.Path
	if len(r.URL.Query()) > 0 {
		text += strings.ReplaceAll(r.URL.Query().Encode(), "%2C", ",")
	}
	log.Println(text + " - " + "[" + string(rune(code)) + "]")
	err := file.Close()
	ErrorHandler(err)
	return
}

func Start(text string) {
	file, issue := os.OpenFile("log/access.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if issue != nil {
		log.Printf("Error opening file: %v", issue)
		return
	}

	log.SetOutput(file)
	log.Println(text)
	err := file.Close()
	ErrorHandler(err)
	return
}