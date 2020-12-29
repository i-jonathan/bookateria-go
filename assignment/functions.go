package assignment

import (
	"bookateriago/core"
	"bookateriago/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDatabase initializes the models for assignments
func initDatabase() *gorm.DB {
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
