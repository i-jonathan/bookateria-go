package document

type Document struct {
	ID		uint 	`json:"id" gorm:"primaryKey;autoIncrement"`
	Title	string 	`json:"title" gorm:"not null;unique"`
	Author	string 	`json:"author"`
	Summary	string 	`json:"summary"`
}
