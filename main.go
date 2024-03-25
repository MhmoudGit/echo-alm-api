package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/MhmoudGit/echo-alm-api/docs"
)

// @title alm-api
// @version 1.0
// @description alm-api server swagger docs

// @host localhost:8000
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
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
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// start server with graceful shutdown
	GracefulShutdown(env, e, &db)
}
