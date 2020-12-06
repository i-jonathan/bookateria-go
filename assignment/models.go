package assignment

import (
	"bookateriago/account"
	"gorm.io/gorm"
	"time"
)

type Problem struct {
	gorm.Model
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Deadline        time.Time    `json:"deadline"`
	User            account.User `json:"user"`
	UserID          int          `json:"user_id"`
	Slug            string       `json:"slug"`
	SubmissionCount int          `json:"submission_count"`
}

type Submission struct {
	gorm.Model
	Problem     Problem      `json:"problem"`
	ProblemID   int          `json:"problem_id"`
	User        account.User `json:"user"`
	UserID      int          `json:"user_id"`
	FileSlug    string       `json:"file_slug"`
	Slug        string       `json:"slug"`
	Submissions int64        `json:"submissions"`
}
