package document

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
)


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

func CheckDuplicate(document *Document) bool {
	var count int64
	db.Model(&Document{}).Where("title = ? AND edition = ? AND author = ?", document.Title, document.Edition, document.Author).Count(&count)	
	return count > 0
}

func FilterBy(queryType string, queryValue string) (bool, []Document) {
	var count int64
	switch queryType {
		case "id":
			id, _ := strconv.Atoi(queryValue)
			db.Model(&Document{}).Where("id = ?", id).Count(&count)
			return count > 0, nil
		case "title":
			//Todo 
		case "author":
			//Todo		
	}
	return true, nil
	
}