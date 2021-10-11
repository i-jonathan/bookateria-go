package forum

import (
	"bookateriago/core"
	"bookateriago/log"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// XExists checks if an object by the slug given exists
//  returns true if it exists, false otherwise
func XExists(slug string, model string) bool {
	var count int64
	var db = InitDatabase()
	switch model {
	case "question":
		db.Model(&question{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "answer":
		db.Model(&answer{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "qUpvote":
		db.Model(&questionUpVote{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "aUpvote":
		db.Model(&answerUpvote{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	default:
		return false
	}
}

// InitDatabase initializes the database and migrates the forum models
func InitDatabase() *gorm.DB {
	var (
		databaseName = os.Getenv("database_name")
		port         = os.Getenv("database_port")
		pass         = os.Getenv("database_pass")
		user         = os.Getenv("database_user")
		host         = os.Getenv("database_host")
		ssl          = os.Getenv("database_ssl")
	)

	postgresConnection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, databaseName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.ErrorHandler(err)

	// err = db.AutoMigrate(&questionTag{}, &oneQuestion{}, &oneAnswer{}, &oneQUpVote{}, &answerUpvote{})
	err = db.AutoMigrate(&questionTag{})
	err = db.AutoMigrate(&question{})
	err = db.AutoMigrate(&answer{})
	err = db.AutoMigrate(&questionUpVote{})
	err = db.AutoMigrate(&answerUpvote{})
	log.ErrorHandler(err)
	return db
}

// validator checks if a string is empty
func validator(values []string) bool {
	for _, value := range values {
		if strings.Join(strings.Fields(value), " ") == "" {
			return false
		}
	}
	return true
}
