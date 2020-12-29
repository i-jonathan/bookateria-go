package forum

import "time"

// QuestionRequest for interfacing with the question request
type QuestionRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// QuestionResponse is the structure of a response returned
type QuestionResponse struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime:nano"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	QuestionTags []TagsResponse
	UpVotes      int `json:"up_votes"`
	UpVoters     []QUpVotesResponse
}

// TagsResponse for returning tags on request
type TagsResponse struct {
	ID   uint
	Name string
}

// QUpVotesResponse responding to question upvotes
type QUpVotesResponse struct {
	User User
}

// User for returning the User details
type User struct {
	UserName  string
	Email     string
	FirstName string
	LastName  string
}
