package document

import (
	"bookateriago/core"
	"bookateriago/log"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
func validate(field string) (string, error) {

	field = strings.ToLower(strings.TrimSpace(field))
	//author := strings.TrimSpace(fields["author"])

	/*title = strings.ToLower(title)
	author = strings.ToLower(author)*/

	if field != "" {
		field = strings.Join(strings.Fields(field), " ")
		return strings.Title(field), nil
	}

	return "", errors.New("Either Title Or Author Is Empty")

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
