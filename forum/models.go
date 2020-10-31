package forum

import (
	"bookateria-api-go/account"
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type QuestionTag struct {
	gorm.Model
	QuestionID uint    `json:"question_id"`
	Name       string `json:"name"`
}

type Question struct {
	gorm.Model
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	QuestionTags []QuestionTag `json:"tags"`
	UserID       int           `json:"user_id"`
	User         account.User  `json:"user"`
	UpVotes      int           `json:"up_votes"`
}

type Answer struct {
	gorm.Model
	QuestionID int          `json:"question_id"`
	Question   Question     `json:"question"`
	Response   string       `json:"response"`
	UpVotes    string       `json:"up_votes"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:SET NULL"`
}

type QuestionUpVote struct {
	gorm.Model
	QuestionID int          `json:"question_id"`
	Question   Question     `json:"question"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
}

type AnswerUpvote struct {
	gorm.Model
	AnswerID int          `json:"answer_id"`
	Answer   Answer       `json:"answer"`
	UserID   int          `json:"user_id"`
	User     account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
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
