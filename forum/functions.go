package forum

import (
	"bookateriago/core"
	"bookateriago/log"
	"fmt"
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
		db.Model(&Question{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "answer":
		db.Model(&Answer{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "qUpvote":
		db.Model(&QuestionUpVote{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "aUpvote":
		db.Model(&AnswerUpvote{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	default:
		return false
	}
}

// InitDatabase initializes the database and migrates the forum models
func InitDatabase() *gorm.DB {
	viperConfig := core.ReadViper()
	var (
		databaseName = viperConfig.Get("database.name")
		port         = viperConfig.Get("database.port")
		pass         = viperConfig.Get("database.pass")
		user         = viperConfig.Get("database.user")
		host         = viperConfig.Get("database.host")
		ssl          = viperConfig.Get("database.ssl")
	)

	postgresConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, databaseName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.ErrorHandler(err)

	// err = db.AutoMigrate(&QuestionTag{}, &Question{}, &Answer{}, &QuestionUpVote{}, &AnswerUpvote{})
	err = db.AutoMigrate(&QuestionTag{})
	err = db.AutoMigrate(&Question{})
	err = db.AutoMigrate(&Answer{})
	err = db.AutoMigrate(&QuestionUpVote{})
	err = db.AutoMigrate(&AnswerUpvote{})
	log.ErrorHandler(err)
	return db
}
