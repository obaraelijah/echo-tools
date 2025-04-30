package middleware

import (
	"github.com/labstack/echo/v4"
)

type AllowedHost struct {
	Host  string
	Https bool
}
type SecurityConfig struct {
	AllowedHosts            []AllowedHost
	UseForwardedProtoHeader bool
}

func Security(config *SecurityConfig) echo.MiddlewareFunc {
	if config == nil {
		panic("Security config must not be nil")
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			allowed := false
			if config.UseForwardedProtoHeader {
				for _, allowedHost := range config.AllowedHosts {
					if c.Request().Host == allowedHost.Host {
						if c.Request().TLS != nil && allowedHost.Https {
							allowed = true
							break
						} else if c.Request().TLS == nil {
							proto := c.Request().Header.Get("X-Forwarded-Proto")
							if !allowedHost.Https {
								if proto == "http" {
									allowed = true
									break
								} else if proto == "" {
									allowed = true
									break
								}
							} else {
								if proto == "https" {
									allowed = true
									break
								}
							}
						}
					}
				}
			} else {
				for _, allowedHost := range config.AllowedHosts {
					if c.Request().Host == allowedHost.Host {
						if c.Request().TLS != nil && allowedHost.Https {
							allowed = true
							break
						} else if c.Request().TLS == nil && !allowedHost.Https {
							allowed = true
							break
						}
					}
				}
			}
			if !allowed {
				proto := "https://"
				if c.Request().TLS == nil {
					proto = "http://"
				}
				c.Logger().Debugf("%s is not in allowed hosts", proto+c.Request().Host)
				return c.JSON(401, struct{ Error string }{Error: "not allowed"})
			}
			return next(c)
		}
	}
}
