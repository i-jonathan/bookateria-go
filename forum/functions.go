package forum

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	log.Handler("panic", "Couldn't connect to DB", err)

	err = db.AutoMigrate(&QuestionTag{})
	err = db.AutoMigrate(&Question{})
	err = db.AutoMigrate(&Answer{})
	err = db.AutoMigrate(&QuestionUpVote{})
	err = db.AutoMigrate(&AnswerUpvote{})
	log.Handler("warn", "Couldn't Migrate model to DB", err)
	return db
}