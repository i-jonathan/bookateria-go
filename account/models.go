package account

import "time"

// User model. Simple enough
type User struct {
	// For Returning Data, might have to create another struct that is used solely for reading from
	// Seems there's no write only for json or gorm for that matter
	ID              uint      `json:"id"`
	UserName        string    `json:"user_name"`
	FullName        string    `json:"full_name" gorm:"not null"`
	Alias	        string    `json:"alias"`
	Email           string    `json:"email" gorm:"not null;unique"`
	IsAdmin         bool      `json:"is_admin" gorm:"default:false"`
	Password        string    `json:"password"`
	LastLogin       time.Time `json:"last_login"`
	IsActive        bool      `json:"is_active" gorm:"default:false"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// profile model. Could be extended soon
type profile struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//Bio		string	`json:"bio"`
	//Picture	string	`json:"picture"`
	Points int  `json:"points" gorm:"default:20"`
	UserID int  `json:"user_id"`
	User   User `json:"user" gorm:"constraints:OnDelete:CASCADE;not null;unique"`
}

// passwordConfig for constructing really nice and secure passwords. Make sense? No
type passwordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}
