package forum

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os/user"
)

type QuestionTag struct {
	gorm.Model
	Name	string	`json:"name"`
}

type Question struct {
	gorm.Model
	Title			string			`json:"title"`
	Description		string			`json:"description"`
	Tags			[]QuestionTag	`json:"tags"`
	AcceptedAnswer	string			`json:"accepted_answer"`
	User			user.User		`json:"user"`
	UpVotes			int				`json:"up_votes"`
}

type Answer struct {
	gorm.Model
	Question	Question	`json:"question"`
	Response	string		`json:"response"`
	UpVotes		string		`json:"up_votes"`
	User		user.User	`json:"user" gorm:"constraints:OnDelete:SET NULL"`
}

type QuestionUpVote struct {
	gorm.Model
	Question	Question	`json:"question"`
	User 		user.User	`json:"user" gorm:"constraints:OnDelete:CASCADE"`
}

type AnswerUpvote struct {
	gorm.Model
	Answer	Answer 		`json:"answer"`
	User	user.User	`json:"user" gorm:"constraints:OnDelete:CASCADE"`
}

func InitDatabase() *gorm.DB {
	viperConfig := core.ReadViper()
	var (
		databaseName = viperConfig.Get("database.name")
		port = viperConfig.Get("database.port")
		pass = viperConfig.Get("database.pass")
		user = viperConfig.Get("database.user")
		host = viperConfig.Get("database.host")
		ssl = viperConfig.Get("database.ssl")
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