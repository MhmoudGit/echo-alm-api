package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

func GracefulShutdown(env ENV, e *echo.Echo, db Database) {
	// server configurations [ starting - graceful shutdown]
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		err := e.Start(env.Address)
		if err != nil && err != http.ErrServerClosed {
			db.disconnect()
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 2 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := e.Shutdown(ctx)
	if err != nil {
		e.Logger.Fatal(err)
	}
	db.disconnect()
	<-ctx.Done()
}
