package assignment

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	err = db.AutoMigrate(&Problem{}, &Submission{})
	log.ErrorHandler(err)

	return db
}

func XExists(slug, model string) bool {
	var count int64
	var db = InitDatabase()

	switch model {
	case "question":
		db.Model(&Problem{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	case "submission":
		db.Model(&Submission{}).Where("slug = ?", slug).Count(&count)
		return count > 0
	default:
		return false
	}
}
