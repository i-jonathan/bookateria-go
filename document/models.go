package document

import"gorm.io/gorm"

type Document struct {
	gorm.Model
	ID		uint 	`json:"id" gorm:"primaryKey;autoIncrement"`
	Title	string 	`json:"title" gorm:"not null"`
	Edition int32 	`json: "edition"`
	Author	string 	`json:"author"`
	Summary	string 	`json:"summary"`
}


