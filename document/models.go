package document

import (
	"bookateria-api-go/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Title   string `json:"title" gorm:"not null;unique"`
	Author  string `json:"author"`
	Summary string `json:"summary"`
}

func InitDatabase() *gorm.DB {
	postgresConnection := "host=localhost port=5432 user=postgres dbname=bookateria-go password=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.Handler("panic", "Couldn't connect to DB", err)

	err = db.AutoMigrate(&Document{})
	log.Handler("warn", "Couldn't Migrate model to DB", err)
	return db
}
