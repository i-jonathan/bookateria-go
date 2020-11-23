package assignment

import (
	"bookateria-api-go/account"
	"gorm.io/gorm"
	"time"
)

type Question struct {
	gorm.Model
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Deadline    time.Time    `json:"deadline"`
	User        account.User `json:"user"`
	Slug        string       `json:"slug"`
}

type Submission struct {
	gorm.Model
	Question Question     `json:"question"`
	User     account.User `json:"user"`
	FileSlug string       `json:"file_slug"`
}
