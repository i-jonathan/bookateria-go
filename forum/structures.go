package forum

import (
	"gorm.io/gorm"
)

type QuestionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type QuestionResponse struct {
	gorm.Model
	Title        string `json:"title"`
	Description  string `json:"description"`
	QuestionTags []TagsResponse
	UpVotes      int `json:"up_votes"`
	UpVoters     []QUpVotesResponse
}

type TagsResponse struct {
	ID   uint
	Name string
}

type QUpVotesResponse struct {
	User User
}

type User struct {
	UserName  string
	Email     string
	FirstName string
	LastName  string
}
