package middleware

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/utilitymodels"
	"gorm.io/gorm"
)

var (
	ErrDatabaseError         = errors.New("there was a problem updating the database")
	ErrCookieNotFound        = errors.New("cookie is missing")
	ErrSessionContextMissing = errors.New("session context is missing")
)

// GetSessionContext returns a SessionContext from a Context
func GetSessionContext(c echo.Context) (SessionContext, error) {
	if context, ok := c.Get("SessionContext").(SessionContext); !ok {
		return nil, ErrSessionContextMissing
	} else {
		return context, nil
	}

}

// Login This method is used to log a user in. auth.Authenticate has to be called before.
// A cookie is set if the user can be logged in.
// Parameter user: Can be retrieved by auth.Authenticate.
// Parameter c: Pointer to the current context. Must implement middleware.SessionContext
// Parameter config: Refer to SessionConfig.
func Login(db *gorm.DB, user *utilitymodels.User, c echo.Context) error {
	context := c.Get("SessionContext").(SessionContext)

	// Couldn't find session with the current user associated
	session := utilitymodels.Session{
		UserID:     user.ID,
		ValidUntil: time.Now().UTC().Add(*context.GetSessionConfig().CookieAge),
	}

	// Generation of session id
	var count int64
	r := make([]byte, 64)
	for {
		if _, err := rand.Read(r); err != nil {
			c.Logger().Warn("Error while generating random numbers")
			continue
		}
		sessionID := fmt.Sprintf("%x", r)
		db.Find(&utilitymodels.Session{}).Where("session_id = ?", sessionID).Count(&count)
		if count == 0 {
			session.SessionID = sessionID
			break
		}
		c.Logger().Debugf("Generated session_id already in database, regenerating ..")
	}

	if err := db.Create(&session).Error; err != nil {
		c.Logger().Errorf("Error saving session to database: %s", err.Error())
		return ErrDatabaseError
	} else {
		if err := db.Model(&user).Update("last_login_at", time.Now().UTC()).Error; err != nil {
			c.Logger().Warnf("Error updating last_login_at of user %d: %s", user.ID, err.Error())
			return ErrDatabaseError
		}

		// Session was saved, we can set the cookie
		cookie := &http.Cookie{
			Name:     context.GetSessionConfig().CookieName,
			Value:    session.SessionID,
			Path:     context.GetSessionConfig().CookiePath,
			Domain:   "", // Only allow current site
			MaxAge:   int(context.GetSessionConfig().CookieAge.Seconds()),
			Secure:   *context.GetSessionConfig().Secure,
			SameSite: http.SameSiteDefaultMode,
		}

		c.SetCookie(cookie)
	}
	return nil
}

// Logout Helper method to logout and therefore invalidating a user's session. If the user isn't logged in,
// nil is returned
func Logout(db *gorm.DB, c echo.Context) error {
	sessionContext := c.Get("SessionContext").(SessionContext)

	// If user is not authenticated, there's nothing to do
	if !sessionContext.IsAuthenticated() {
		return ErrCookieNotFound
	}

	if err := db.Where("session_id = ?", *sessionContext.GetSessionID()).Delete(&utilitymodels.Session{}).Error; err != nil {
		c.Logger().Error(err.Error())
		return ErrDatabaseError
	}

	c.SetCookie(&http.Cookie{
		Name:   sessionContext.GetSessionConfig().CookieName,
		Value:  "",
		Path:   sessionContext.GetSessionConfig().CookiePath,
		Domain: "",
		MaxAge: -1, // Cookie is invalidated through MaxAge < 0
		Secure: *sessionContext.GetSessionConfig().Secure,
	})

	// Flushing current session
	sessionContext.flush()
	return nil
}

// InvalidateSessions Helper method to invalidate all sessions of a user
func InvalidateSessions(db *gorm.DB, userID uint) error {
	if err := db.Where("user_id = ?", userID).Delete(&utilitymodels.Session{}).Error; err != nil {
		return ErrDatabaseError
	}
	return nil
}
