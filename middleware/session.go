package middleware

import (
	"errors"
	"reflect"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/logging"
	"github.com/obaraelijah/echo-tools/utilitymodels"
	"gorm.io/gorm"
)

var ErrSessionMisconfigured = errors.New(
	"field SessionContext does not exist in Context. Skipping this middleware. " +
		"Check if CustomContext is enabled and SessionContext is embedded in your context struct",
)

type SessionContext interface {
	GetUserID() *uint
	GetSessionID() *string
	IsAuthenticated() bool
	GetSessionConfig() *SessionConfig
	flush()
}

// SessionConfig Set the parameters for the Session.
// Parameter CookieName defaults to "session_id".
// Parameter CookieAge defaults to 30 * time.Minute.
// Parameter Secure defaults to true. If set, the cookie can only be sent through an HTTPS connection.
// Parameter CookiePath defaults to "". Can be used to restrict the path the cookie can be sent to.
// Parameter DisableLogging defaults to false. If set, no debug logs are sent. Error logs are still sent.
type SessionConfig struct {
	CookieName     string
	CookieAge      *time.Duration
	Secure         *bool
	CookiePath     string
	DisableLogging bool
}

type s struct {
	userID        *uint
	authenticated bool
	sessionConfig *SessionConfig
	sessionID     *string
}

// GetUserID Returns a pointer to a userID if present else nil
func (s *s) GetUserID() *uint {
	return s.userID
}

// IsAuthenticated Returns true if the session of this request is valid
func (s *s) IsAuthenticated() bool {
	return s.authenticated
}

// GetSessionConfig Returns the config of the session middleware. Mostly for internal use of session.Login
func (s *s) GetSessionConfig() *SessionConfig {
	return s.sessionConfig
}

func (s *s) GetSessionID() *string {
	return s.sessionID
}

func (s *s) flush() {
	s.userID = nil
	s.authenticated = false
	s.sessionID = nil
}

func (config *SessionConfig) FixSessionConfig() {
	if config.CookieName == "" {
		config.CookieName = "session_id"
	}
	if config.Secure == nil {
		secure := true
		config.Secure = &secure
	}
	if config.CookieAge == nil {
		age := 30 * time.Minute
		config.CookieAge = &age
	}
	return
}

// Session Use as middleware. Requires CustomContext to be set with a corresponding struct that embeds SessionContext
// or has a field named SessionContext. If SessionContext is not found, the middleware is skipped.
func Session(db *gorm.DB, log logging.Logger, config *SessionConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &SessionConfig{}
	}
	config.FixSessionConfig()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if SessionContext is available
			field := reflect.ValueOf(c).Elem().FieldByName("SessionContext")
			if field == (reflect.Value{}) {
				log.Error(ErrSessionMisconfigured.Error())
				// Skipping middleware to not break the server
				return next(c)
			}

			sessionContext := &s{
				userID:        nil,
				authenticated: false,
				sessionConfig: config,
			}

			// Check if cookie is present
			if cookie, err := c.Cookie(config.CookieName); err != nil {
				// No need to do something, default values of sessionContext are fine
				if !config.DisableLogging {
					log.Debugf("Cookie \"%s\" is not present in request", config.CookieName)
				}
			} else {

				var sessionCount int64
				var session utilitymodels.Session
				db.Find(&session).Where("session_id = ?", cookie.Value).Count(&sessionCount)
				switch sessionCount {
				case 0:
					// No session with that id was found
					if !config.DisableLogging {
						log.Debugf("Cookie with SessionID %s was not found in DB", cookie.Value)
					}
				case 1:
					// Session was found

					// Check if session is not expired
					if !time.Now().UTC().After(session.ValidUntil) {
						var user utilitymodels.User
						if db.Model(session).Association("User").Find(&user); err != nil {
							log.Warn(err.Error())
						} else {
							// Check if user is valid
							if user.ID > 0 && user.Active.Valid && user.Active.Bool {
								sessionContext.userID = &user.ID
								sessionContext.sessionID = &session.SessionID
								sessionContext.authenticated = true
							} else {
								// User is invalid or not active
								if !config.DisableLogging {
									log.Debugf(
										"Invalid or deactivated user: userID: %d | %+v",
										user.ID, user.Active,
									)
								}
							}
						}
					}
				default:
					// This is broken!1!!1elf!
				}
			}

			// Set SessionContext
			field.Set(reflect.ValueOf(&sessionContext).Elem())
			return next(c)
		}
	}
}
