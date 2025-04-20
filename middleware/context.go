package middleware

import (
	"reflect"

	"github.com/labstack/echo/v4"
)

// Wrap helper function to pass functions using CustomContext to GET, POST, PUT, DELETE, ... of echo directly
func Wrap[T any](f func(cc T) error) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(T)
		return f(cc)
	}
}

// CustomContext Requires a pointer to a struct that embeds echo.Context or has a field with the name Context
func CustomContext(cc any) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reflect.ValueOf(cc).Elem().FieldByName("Context").Set(reflect.ValueOf(&c).Elem())
			return next(cc.(echo.Context))
		}
	}
}
