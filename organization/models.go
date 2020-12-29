package organization

import (
	"bookateriago/account"
	"time"
)

// class -
type class struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"`
	Name      string         `json:"name"`
	Owner     account.User   `json:"owner"`
	Members   []account.User `json:"members"`
}
