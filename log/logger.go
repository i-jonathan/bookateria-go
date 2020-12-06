package log

import (
	"log"
	"os"
)

func ErrorHandler(err error) {
	file, issue := os.OpenFile("log/error.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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

func AccessHandler(text string) {
	file, issue := os.OpenFile("log/access.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
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