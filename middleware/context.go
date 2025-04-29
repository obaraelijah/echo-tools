package middleware

import (
	"github.com/labstack/echo/v4"
)

// Wrap helper function to pass functions using CustomContext to GET, POST, PUT, DELETE, ... of echo directly
func Wrap[T any](f func(cc T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(T)
		return f(cc)
	}
}
