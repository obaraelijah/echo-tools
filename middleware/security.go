package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/logging"
	"github.com/obaraelijah/echo-tools/utility"
)

type AllowedHost struct {
	Host  string
	Https bool
}

type SecurityConfig struct {
	AllowedHosts            []AllowedHost
	UseForwardedProtoHeader bool
}

func Security(log logging.Logger, config *SecurityConfig) echo.MiddlewareFunc {
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
							if !allowedHost.Https {
								allowed = true
								break
							} else {
								proto := c.Request().Header.Get("X-Forwarded-Proto")
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
				log.Debugf("%s is not in allowed hosts", proto+c.Request().Host)
				return c.JSON(401, utility.JsonResponse{Error: "not allowed"})
			}
			return next(c)
		}
	}
}
