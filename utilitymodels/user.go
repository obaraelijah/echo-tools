package utilitymodels

import (
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type LocalUser struct {
	Common
	LastLoginAt sql.NullTime `json:"-" gorm:"default:null"` // This is only relevant if the session middleware is in use
	Email       *string      `json:"email" gorm:"unique;default:null"`
	Username    string       `json:"username" gorm:"unique;not null"`
	Password    string       `json:"-" gorm:"not null"`
}

type LDAPProvider struct {
	Common
	Name         string
	Uri          string
	SearchBase   string
	SearchFilter *string
	AdminGroup   *string
}

type LDAPUser struct {
	Common
	LastLoginAt    sql.NullTime `json:"-" gorm:"default:null"` // This is only relevant if the session middleware is in use
	LDAPProviderID uint
	LDAPProvider   LDAPProvider
	Username       string
}

func (user *LDAPUser) GetAuthModelIdentifier() (string, uint) {
	return "ldap", user.ID
}

func (user *LDAPUser) UpdateLastLogin(echo.Context, *gorm.DB, time.Time) {

}

func GetLDAPUser(db *gorm.DB) func() (string, func(foreignKey uint) any) {
	return func() (string, func(foreignKey uint) any) {
		return "ldap", func(foreignKey uint) any {
			var user LDAPUser

			var count int64
			db.Find(&user, "ID = ?", foreignKey).Count(&count)

			if count != 1 {
				return nil
			}

			return &user
		}
	}
}

func (user *LocalUser) GetAuthModelIdentifier() (string, uint) {
	return "local", user.ID
}

func (user *LocalUser) UpdateLastLogin(c echo.Context, db *gorm.DB, loginTime time.Time) {
	if err := db.Model(&user).Update("last_login_at", loginTime).Error; err != nil {
		c.Logger().Warnf("Error updating last_login_at of user %d: %s", user.ID, err.Error())
	}
}

func GetLocalUser(db *gorm.DB) func() (string, func(foreignKey uint) any) {
	return func() (string, func(foreignKey uint) any) {
		return "local", func(foreignKey uint) any {
			var user LocalUser

			var count int64
			db.Find(&user, "ID = ?", foreignKey).Count(&count)

			if count != 1 {
				return nil
			}

			return &user
		}
	}
}
