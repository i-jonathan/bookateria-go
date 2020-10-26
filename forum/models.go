package forum

import (
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

func InitDatabase(m)  {
	
}