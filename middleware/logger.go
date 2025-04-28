package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/logging"
)

func Logging(log logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			if err := next(c); err != nil {
				switch err.(type) {
				case *echo.HTTPError:
					httpError := err.(*echo.HTTPError)
					log.Infof(
						"%d %s %s %v - %s %s",
						httpError.Code, c.Request().Method, c.Request().RequestURI, time.Now().Sub(start),
						c.RealIP(), c.Request().Header.Get("User-Agent"),
					)
				default:
					log.Error(err.Error())
				}
			} else {
				log.Infof(
					"%d %s %s %v - %s %s",
					c.Response().Status, c.Request().Method, c.Request().RequestURI, time.Now().Sub(start),
					c.RealIP(), c.Request().Header.Get("User-Agent"),
				)
			}
			return nil
		}
	}
}
