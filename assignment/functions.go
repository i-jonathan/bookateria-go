package assignment

import (
	"bookateriago/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// InitDatabase initializes the models for assignments
func initDatabase() *gorm.DB {
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

	err = db.AutoMigrate(&problem{}, &submission{})
	log.ErrorHandler(err)

	return db
}

// XExists checks the existence of an object given the slug and the model
func xExists(slug, model string) bool {
	var count int64
	var db = initDatabase()

	switch model {
	case "question":
		db.Model(&problem{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "submission":
		db.Model(&submission{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	default:
		return false
	}
}
