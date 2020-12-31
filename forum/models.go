package forum

import (
	"bookateriago/account"
	"time"
)

// questionTag model for tags attached to questions
type questionTag struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID uint      `json:"question_id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
}

// question is the model for forum questions
type question struct {
	ID           uint          `json:"id"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime:nano"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	QuestionTags []questionTag `json:"tags" gorm:"foreignKey:QuestionID"`
	UserID       int           `json:"user_id"`
	User         account.User  `json:"user"`
	//UpVotes      int           `json:"up_votes"`
	UpVoters []questionUpVote `json:"up_voters" gorm:"foreignKey:QuestionID"`
	Slug     string           `json:"slug"`
}

// answer are replies to Questions
type answer struct {
	ID         uint         `json:"id"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID int          `json:"question_id"`
	Question   question     `json:"question"`
	Response   string       `json:"response"`
	UpVotes    string       `json:"up_votes"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:SET NULL"`
	Slug       string       `json:"slug"`
}

// questionUpVote for keeping a list of up votes on a oneQuestion
type questionUpVote struct {
	ID         uint         `json:"id"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID int          `json:"question_id"`
	Question   question     `json:"question"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
}

// answerUpvote for keeping a lost of upvotes on an answer
type answerUpvote struct {
	ID        uint         `json:"id"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	AnswerID  int          `json:"answer_id"`
	Answer    answer       `json:"answer"`
	UserID    int          `json:"user_id"`
	User      account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
}
