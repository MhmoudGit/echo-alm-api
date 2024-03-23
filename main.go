package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load the environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file")
	}
	env := SetEnv()

	// connecting database
	var db Storage
	db.GormConnect(env)
	db.GormAutoMigrateDb()

	e := echo.New()
	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello World!")
	})

	// start server with graceful shutdown
	GracefulShutdown(env, e, &db)
}

func GracefulShutdown(env ENV, e *echo.Echo, db *Storage) {
	// server configurations [ starting - graceful shutdown]
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		err := e.Start(env.Address)
		if err != nil && err != http.ErrServerClosed {
			db.GormClose()
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
	db.GormClose()
	<-ctx.Done()
}
