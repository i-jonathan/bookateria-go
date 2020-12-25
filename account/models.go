package account

import (
	"gorm.io/gorm"
	"time"
)

// User model. Simple enough
type User struct {
	// For Returning Data, might have to create another struct that is used solely for reading from
	// Seems there's no write only for json or gorm for that matter
	gorm.Model
	UserName        string    `json:"user_name"`
	FirstName       string    `json:"first_name" gorm:"not null"`
	LastName        string    `json:"last_name" gorm:"not null"`
	Email           string    `json:"email" gorm:"not null;unique"`
	IsAdmin         bool      `json:"is_admin" gorm:"default:false"`
	Password        string    `json:"password"`
	LastLogin       time.Time `json:"last_login"`
	IsActive        bool      `json:"is_active" gorm:"default:false"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"default:false"`
}

// Profile model. Could be extended soon
type Profile struct {
	gorm.Model
	//Bio		string	`json:"bio"`
	//Picture	string	`json:"picture"`
	Points int  `json:"points" gorm:"default:20"`
	UserID int  `json:"user_id"`
	User   User `json:"user" gorm:"constraints:OnDelete:CASCADE;not null;unique"`
}

// PasswordConfig for contructing really nice and secure passwords. Make sense? No
type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}
