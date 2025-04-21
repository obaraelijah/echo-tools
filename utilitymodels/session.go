package utilitymodels

import "time"

type Session struct {
	Common
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user" gorm:"not null;constraint:OnDelete:CASCADE"`
	SessionID  string    `json:"-" gorm:"not null;unique"`
	ValidUntil time.Time `json:"valid_until" gorm:"not null"`
}
