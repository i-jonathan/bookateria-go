package forum

import (
	"bookateriago/account"
	"gorm.io/gorm"
)

type QuestionTag struct {
	gorm.Model
	QuestionID uint   `json:"question_id"`
	Name       string `json:"name"`
	Slug       string `json:"slug"`
}

type Question struct {
	gorm.Model
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	QuestionTags []QuestionTag `json:"tags"`
	UserID       int           `json:"user_id"`
	User         account.User  `json:"user"`
	//UpVotes      int           `json:"up_votes"`
	UpVoters []QuestionUpVote `json:"up_voters"`
	Slug     string           `json:"slug"`
}

type Answer struct {
	gorm.Model
	QuestionID int          `json:"question_id"`
	Question   Question     `json:"question"`
	Response   string       `json:"response"`
	UpVotes    string       `json:"up_votes"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:SET NULL"`
	Slug       string       `json:"slug"`
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
