package utilitymodels

import "time"

type Session struct {
	Common
	AuthID     uint      `json:"auth_id" gorm:"not null"`
	AuthKey    string    `json:"auth_key" gorm:"not null"`
	SessionID  string    `json:"-" gorm:"not null;unique"`
	ValidUntil time.Time `json:"valid_until" gorm:"not null"`
}
