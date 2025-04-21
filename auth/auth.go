package auth

import (
	"errors"
	"fmt"

	"github.com/obaraelijah/echo-tools/db"
	"github.com/obaraelijah/echo-tools/utilitymodels"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrUsernameNotFound     = errors.New("username not found")
)

// Authenticate Try to authenticate with the given credentials
func Authenticate(username string, password string) (*utilitymodels.User, error) {
	var u utilitymodels.User
	var count int64

	db.DB.Find(&u, "username = ?", username).Count(&count)
	if count == 0 {
		// Comparing static hash in order to deny username enumeration by looking at the time a request took
		bcrypt.CompareHashAndPassword(
			[]byte("$2b$12$KisigGoquLISbypB3kHB1eUOXZUWm7AwOZcwIIH9V9YejhxvIvlo6"),
			[]byte("Deny username enumeration"),
		)
		return nil, ErrUsernameNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		fmt.Println(err.Error())
		return nil, ErrAuthenticationFailed
	}

	return &u, nil
}
