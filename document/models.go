package document

import"gorm.io/gorm"

type Tag struct {
	DocumentID uint `json:"documentid"`
	TagName string `json:"tagname"`
}

type Document struct {
	gorm.Model
	ID		uint 	`json:"id" gorm:"primaryKey;autoIncrement"`
	Title	string 	`json:"title" gorm:"not null"`
	Edition int32 	`json: "edition"`
	Author	string 	`json:"author"`
	Summary	string 	`json:"summary"`
	Tags	[]Tag 	`json:"tags"`
}


