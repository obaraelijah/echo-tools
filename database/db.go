package database

import (
	"os"

	"github.com/obaraelijah/echo-tools/utilitymodels"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Initialize(dial gorm.Dialector, models ...interface{}) *gorm.DB {
	// Open DB
	conn, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}

	models = append(models, &utilitymodels.LocalUser{})
	models = append(models, &utilitymodels.LDAPUser{})
	models = append(models, &utilitymodels.LDAPProvider{})
	models = append(models, &utilitymodels.Session{})

	// Migrate
	if err := conn.AutoMigrate(
		models...,
	); err != nil {
		os.Exit(1)
	}

	return conn
}

// CreateLocalUser Helper method to create a user. bcrypt with a cost of 12 is used as hash.
func CreateLocalUser(db *gorm.DB, username string, password string, email *string) (*utilitymodels.LocalUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	u := utilitymodels.LocalUser{
		Username: username,
		Email:    email,
		Password: string(hash),
	}
	if err := db.Create(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}
