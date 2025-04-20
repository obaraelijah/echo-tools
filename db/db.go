package db

import (
	"github.com/obaraelijah/echo-tools/utilitymodels"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(dial gorm.Dialector, models ...interface{}) {
	// Open DB
	conn, err := gorm.Open(dial, &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	// Migrate
	if err := conn.AutoMigrate(
		&utilitymodels.User{},
		models,
	); err != nil {
		panic(err.Error())
	}

	DB = conn
}

func CreateUser(username string, password string, email string) (*utilitymodels.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	u := utilitymodels.User{
		Username: username,
		Email:    email,
		Password: string(pw),
	}

	if err := DB.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
