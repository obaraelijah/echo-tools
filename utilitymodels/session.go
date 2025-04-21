package utilitymodels

import "time"

type Session struct {
	CommonProps
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user" gorm:"not null"`
	SessionID  string    `json:"-" gorm:"not null;unique"`
	ValidUntil time.Time `json:"valid_until" gorm:"not null"`
}
