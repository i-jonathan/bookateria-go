package document

import (
	"bookateriago/core"
	"bookateriago/log"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

func InitDatabase() *gorm.DB {
	viperConfig := core.ReadViper()

	var (
		dbName = viperConfig.Get("database.name")
		port   = viperConfig.Get("database.port")
		pass   = viperConfig.Get("database.pass")
		user   = viperConfig.Get("database.user")
		host   = viperConfig.Get("database.host")
		ssl    = viperConfig.Get("database.ssl")
	)
	postgresConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.ErrorHandler(err)

	err = db.AutoMigrate(&Document{})
	err = db.AutoMigrate(&Tag{})

	log.ErrorHandler(err)
	return db
}

func checkDuplicate(document *Document) bool {
	var count int64
	db.Model(&Document{}).Where("title LIKE ? AND edition = ? AND author LIKE ? ", document.Title, document.Edition, document.Author).Count(&count)
	return count > 0
}

func xExists(id uint) bool {
	var count int64
	db.Model(&Document{}).Where("id = ?", id).Count(&count)
	return count > 0
}

/*validate=========================
params: r type: http.Request
Returns: string, string, int, error
===================================*/
func validate(r *http.Request) (string, string, int, error) {

	title := strings.TrimSpace(r.FormValue("title"))
	author := strings.TrimSpace(r.FormValue("author"))
	edition, err := strconv.Atoi(r.FormValue("edition"))

	title = strings.ToLower(title)
	author = strings.ToLower(author)

	if err != nil {
		return "", "", -1, errors.New("Value of edition is not of valid type")
	}

	if title != "" && author != "" {
		title = strings.Join(strings.Fields(title), " ")
		author = strings.Join(strings.Fields(author), " ")
		return strings.Title(title), strings.Title(author), edition, err
	}

	return "", "", -1, errors.New("Either Title Or Author Is Empty")

}

func search(queryType string, queryValue string) []Document {
	switch queryType {
	case "title":
		//Todo
	case "author":
		//Todo
	}
	return nil

}
