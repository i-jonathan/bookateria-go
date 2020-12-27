package organization

import (
	"bookateriago/account"
	"time"
)

// Class -
type Class struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime:nano"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime:nano"`
	Name      string         `json:"name"`
	Owner     account.User   `json:"owner"`
	Members   []account.User `json:"members"`
}
