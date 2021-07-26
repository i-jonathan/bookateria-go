package assignment

import (
	"bookateriago/account"
	"time"
)

// problem is the model for creating assignment questions or problems
type problem struct {
	ID              uint         `json:"id"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Title           string       `json:"title"`
	Description     string       `json:"description"`
	Deadline        time.Time    `json:"deadline"`
	User            account.User `json:"user"`
	UserID          int          `json:"user_id"`
	Slug            string       `json:"slug"`
	SubmissionCount int          `json:"submission_count"`
}

// submission is the model for storing submission data
type submission struct {
	ID          uint         `json:"id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Problem     problem      `json:"problem"`
	ProblemID   int          `json:"problem_id"`
	User        account.User `json:"user"`
	UserID      int          `json:"user_id"`
	FileSlug    string       `json:"file_slug"`
	Slug        string       `json:"slug"`
	Submissions int64        `json:"submissions"`
}
