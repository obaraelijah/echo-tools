package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/utilitymodels"
	"gorm.io/gorm"
)

type SessionContext interface {
	GetUser() any
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
	authModelID   uint
	authModelKey  string
	authenticated bool
	sessionConfig *SessionConfig
	sessionID     *string
}

// GetUser Returns the pointer to a user model or nil if the request was unauthenticated
func (s *s) GetUser() any {
	if f, exists := authMap[s.authModelKey]; !exists {
		return nil
	} else {
		return f(s.authModelID)
	}
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
	s.authModelKey = ""
	s.authModelID = 0
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
	if config.CookiePath == "" {
		config.CookiePath = "/"
	}
	return
}

var authMap map[string]func(uint) any

// RegisterAuthProvider is used to register a new auth provider besides the already existing ones.
// The authIdentifier must be unique and is used to retrieve the correct UserModel with getUserModel which
// should return its own user struct
func RegisterAuthProvider(getProviderInformation func() (string, func(uint) any)) {
	authIdentifier, getUserModel := getProviderInformation()

	if authIdentifier == "" {
		panic("invalid auth provider identifier")
	}

	if _, exists := authMap[authIdentifier]; exists {
		panic("auth provider with that key already exists")
	}

	authMap[authIdentifier] = getUserModel
}

// Session Use as middleware. Requires CustomContext to be set with a corresponding struct that embeds SessionContext
// or has a field named SessionContext. If SessionContext is not found, the middleware is skipped.
func Session(db *gorm.DB, config *SessionConfig) echo.MiddlewareFunc {
	if config == nil {
		config = &SessionConfig{}
	}
	config.FixSessionConfig()

	authMap = map[string]func(foreign uint) any{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check if SessionContext is available
			sessionContext := &s{
				authModelID:   0,
				authModelKey:  "",
				authenticated: false,
				sessionConfig: config,
			}

			// Check if cookie is present
			if cookie, err := c.Cookie(config.CookieName); err != nil {
				// No need to do something, default values of sessionContext are fine
				if !config.DisableLogging {
					c.Logger().Debugf("Cookie \"%s\" is not present in request", config.CookieName)
				}
			} else {

				var sessionCount int64
				var session utilitymodels.Session
				db.Find(&session).Where("session_id = ?", cookie.Value).Count(&sessionCount)
				switch sessionCount {
				case 0:
					// No session with that id was found
					if !config.DisableLogging {
						c.Logger().Debugf("Cookie with SessionID %s was not found in DB", cookie.Value)
					}
				case 1:
					// Session was found

					// Check if session is not expired
					if !time.Now().UTC().After(session.ValidUntil) {
						sessionContext.authModelKey = session.AuthKey
						sessionContext.authModelID = session.AuthID
						sessionContext.sessionID = &session.SessionID

						if sessionContext.GetUser() != nil {
							sessionContext.authenticated = true
						}
					}

				default:
					// This is broken!1!!1elf!
				}
			}

			// Set SessionContext
			c.Set("SessionContext", sessionContext)
			return next(c)
		}
	}
}
