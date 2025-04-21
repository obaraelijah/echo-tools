package db

import (
	"database/sql"
	"os"

	"github.com/obaraelijah/echo-tools/utilitymodels"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(dial gorm.Dialector, models ...interface{}) {
	// Open DB
	conn, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}

	models = append(models, &utilitymodels.User{})

	// Migrate
	if err := conn.AutoMigrate(
		models...,
	); err != nil {
		os.Exit(1)
	}

	DB = conn
}

// CreateUser Helper method to create a user. bcrypt with a cost of 12 is used as hash.
func CreateUser(username string, password string, email string, active bool) (*utilitymodels.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	u := utilitymodels.User{
		Username: username,
		Email:    email,
		Password: string(hash),
		Active: sql.NullBool{
			Bool:  active,
			Valid: true,
		},
	}
	if err := DB.Create(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}
