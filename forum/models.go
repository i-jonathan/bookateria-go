package forum

import (
	"bookateriago/account"
	"time"
)

// QuestionTag model for tags attached to questions
type QuestionTag struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID uint      `json:"question_id"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
}

// Question is the model for forum questions
type Question struct {
	ID           uint          `json:"id"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime:nano"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	QuestionTags []QuestionTag `json:"tags"`
	UserID       int           `json:"user_id"`
	User         account.User  `json:"user"`
	//UpVotes      int           `json:"up_votes"`
	UpVoters []QuestionUpVote `json:"up_voters"`
	Slug     string           `json:"slug"`
}

// Answer are replies to Questions
type Answer struct {
	ID         uint         `json:"id"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID int          `json:"question_id"`
	Question   Question     `json:"question"`
	Response   string       `json:"response"`
	UpVotes    string       `json:"up_votes"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:SET NULL"`
	Slug       string       `json:"slug"`
}

// QuestionUpVote for keeping a list of up votes on a question
type QuestionUpVote struct {
	ID         uint         `json:"id"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	QuestionID int          `json:"question_id"`
	Question   Question     `json:"question"`
	UserID     int          `json:"user_id"`
	User       account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
}

// AnswerUpvote for keeping a lost of upvotes on an answer
type AnswerUpvote struct {
	ID        uint         `json:"id"`
	CreatedAt time.Time    `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time    `json:"updated_at" gorm:"autoUpdateTime:nano"`
	AnswerID  int          `json:"answer_id"`
	Answer    Answer       `json:"answer"`
	UserID    int          `json:"user_id"`
	User      account.User `json:"user" gorm:"constraints:OnDelete:CASCADE"`
}
