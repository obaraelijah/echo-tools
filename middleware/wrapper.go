package middleware

import (
	"github.com/labstack/echo/v4"
)

// LoginRequired Helper function to mark endpoint as login only. Requires Session as a middleware.
// Returns ErrSessionMisconfigured if middleware.SessionContext is not present in Context.
func LoginRequired(f func(echo.Context) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Check if SessionContext is available
		sessionContext := c.Get("SessionContext").(SessionContext)

		// Check if user is authenticated
		if !sessionContext.IsAuthenticated() {
			return c.JSON(403, struct{ Error string }{Error: "Unauthenticated"})
		}

		return f(c)
	}
}
