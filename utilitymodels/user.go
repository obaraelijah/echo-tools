package utilitymodels

import (
	"time"
)

type User struct {
	CommonProps
	LastLoginAt time.Time `json:"-"` // This is only relevant if the session middleware is in use
	Email       string    `json:"email" gorm:"unique"`
	Username    string    `json:"username" gorm:"unique;not null"`
	Password    string    `json:"-" gorm:"not null"`
}
