package utilitymodels

import (
	"time"
)

type CommonID struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type Common struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
