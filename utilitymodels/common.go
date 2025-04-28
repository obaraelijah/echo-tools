package utilitymodels

import (
	"time"

	"gorm.io/gorm"
)

type CommonID struct {
	ID uint `gorm:"primarykey" json:"id"`
}

type Common struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type CommonSoftDelete struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
