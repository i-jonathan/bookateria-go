package document

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Document struct {
	ID		uint 	`json:"id" gorm:"primaryKey;autoIncrement"`
	Title	string 	`json:"title" gorm:"not null;unique"`
	Author	string 	`json:"author"`
	Summary	string 	`json:"summary"`
}

func InitDatabase() *gorm.DB {
	viperConfig := core.ReadViper()


	var(
		dbName = viperConfig.Get("database.name")
		port = viperConfig.Get("database.port")
		pass = viperConfig.Get("database.pass")
		user = viperConfig.Get("database.user")
		host = viperConfig.Get("database.host")
		ssl = viperConfig.Get("database.ssl")
	)
	postgresConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.Handler("panic", "Couldn't connect to DB\n", err)
	
	err = db.AutoMigrate(&Document{})

	log.Handler("warn", "Couldn't Migrate model to DB\n", err)
	return db
}

