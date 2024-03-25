package main

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/MhmoudGit/echo-alm-api/auth"
	_ "github.com/MhmoudGit/echo-alm-api/docs"
)

//	@title			alm-api
//	@version		1.0
//	@description	alm-api server swagger docs

// @host						localhost:8000
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @BasePath					/
func main() {
	// Load the environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Println("error loading .env file")
	}
	env := SetEnv()

	// connecting postgresql database
	db := &Postgre{}
	databaseInit(db, env)
	db.migrate() // this is for postgresql data only

	e := echo.New()

	// middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = &CustomValidator{Validator: validator.New()}

	// Configure middleware with the custom claims type
	jwtconfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}

	// routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	auth.AuthRoutes(e, jwtconfig, db.gorm)

	// start server with graceful shutdown
	GracefulShutdown(env, e, db)
}

// validations for incoming data from the client
type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(400, err.Error())
	}
	return nil
}
