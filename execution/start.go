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

func SignalStart(e *echo.Echo, listenAddress string, restartFunc func()) {
	control := make(chan os.Signal, 1)
	signal.Notify(control, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		// Start server
		if err := e.Start(listenAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal()
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
			cancel()
			break
		} else if sig == syscall.SIGTERM { // Shutdown immediately
			e.Close()
			color.Println(color.PURPLE, "Server was shut down")
			break
		} else {
			fmt.Printf("Received unknown signal: %s\n", sig.String())
		}
	}

	if restart {
		restartFunc()
	}
}
