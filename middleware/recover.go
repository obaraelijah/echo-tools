package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Panic() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					if r == http.ErrAbortHandler {
						panic(r)
					}
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					stack := make([]byte, 4<<10) // 4 KB Stack size
					var length int
					if log.Level() == log.DEBUG {
						length = runtime.Stack(stack, true)
					} else {
						length = runtime.Stack(stack, false)
					}
					stack = stack[:length]

					log.Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
				}
			}()
			return next(c)
		}
	}
}
