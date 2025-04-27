package execution

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/obaraelijah/echo-tools/color"
)

type Config struct {
	ReloadFunc    func()
	StopFunc      func()
	TerminateFunc func()
}

func SignalStart(e *echo.Echo, listenAddress string, config *Config) {
	control := make(chan os.Signal, 1)
	signal.Notify(control, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		// Start server
		if err := e.Start(listenAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println(err.Error())
		}
	}()

	restart := false
	for {
		sig := <-control

		if sig == syscall.SIGHUP { // Reload server
			color.Println(color.PURPLE, "Server is restarting")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			e.Shutdown(ctx)
			cancel()
			restart = true
			break
		} else if sig == syscall.SIGINT { // Shutdown gracefully
			color.Println(color.PURPLE, "Server is stopping gracefully")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			e.Shutdown(ctx)
			config.StopFunc()
			cancel()
			break
		} else if sig == syscall.SIGTERM { // Shutdown immediately
			e.Close()
			config.TerminateFunc()
			color.Println(color.PURPLE, "Server was shut down")
			break
		} else {
			fmt.Printf("Received unknown signal: %s\n", sig.String())
		}
	}

	if restart {
		config.ReloadFunc()
	}
}
