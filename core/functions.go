package core

import (
	"bookateria-api-go/document"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func ErrHandler(err error) {
	if err != nil {
		log.Println(err)
	}
	return
}

func InitDatabase() *gorm.DB {
	postgresConnection := "host=localhost port=5432 user=postgres dbname=bookateria-go password=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	ErrHandler(err)

	err = db.AutoMigrate(&document.Document{})
	ErrHandler(err)
	return db
}
