package organization

import (
	"gorm.io/gorm"
	"bookateriago/account"
)

// Class -
type Class struct {
	gorm.Model
	Name    string 		`json:"name"`
	Owner 	account.User `json:"owner"`
	Members []account.User `json:"members"`

}
