package main

import (
	"log"

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

	// start server with graceful shutdown
	GracefulShutdown(env, e, &db)
}
