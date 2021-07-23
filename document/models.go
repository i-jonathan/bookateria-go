package document

import "bookateriago/account"
import "time"

type Category struct {
	ID           uint   `json:"id" gorm:"primaryKey;autoIncrement; unique"`
	DocumentID   uint   `json:"document_id"`
	CategoryName string `json:"category_name"`
	Slug         string `json:"category_slug"`
}

type Tag struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement; unique"`
	DocumentID uint   `json:"document_id"`
	TagName    string `json:"tag_name"`
	Slug       string `json:"tag_slug"`
}

type Document struct {
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	Size       float64       `json:"size"`
	Downloads  int64        `json:"downloads"`
	ID         uint         `json:"id" gorm:"primaryKey;autoIncrement; unique"`
	Title      string       `json:"title" gorm:"not null"`
	Edition    int          `json:"edition" gorm:"default:0"`
	Author     string       `json:"author"`
	Summary    string       `json:"summary"`
	Tags       []Tag        `json:"tags"`
	FileSlug   string       `json:"file_slug"`
	Slug       string       `json:"slug"`
	CoverSlug  string       `json:"cover_slug"`
	UploaderID int          `json:"uploader_id"`
	Uploader   account.User `json:"uploader"`
	Category   Category     `json:"category"`
}
