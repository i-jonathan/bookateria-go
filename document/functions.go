package document

import (
	"bookateriago/log"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func InitDatabase() *gorm.DB {

	var (
		dbName = os.Getenv("database_name")
		port   = os.Getenv("database_port")
		pass   = os.Getenv("database_pass")
		user   = os.Getenv("database_user")
		host   = os.Getenv("database_host")
		ssl    = os.Getenv("database_ssl")
	)
	postgresConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.ErrorHandler(err)

	err = db.AutoMigrate(&Document{})
	err = db.AutoMigrate(&Tag{})
	err = db.AutoMigrate(&Category{})

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

	return "", errors.New("either title or author is empty")

}

func Paginate(r *http.Request) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}

		pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		switch {
		case pageSize > 50:
			pageSize = 50
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
