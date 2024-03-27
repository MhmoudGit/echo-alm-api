package auth

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AuthRoutes(e *echo.Echo, config echojwt.Config, db *gorm.DB, secret string, sender EmailSender) {
	auth := e.Group("/auth")

	auth.POST("/token", func(c echo.Context) error { return tokenHandler(c, db, secret) })
	auth.POST("/refresh-token", func(c echo.Context) error { return refreshTokenHandler(c, secret) })
	auth.POST("/signup", func(c echo.Context) error { return signupHandler(c, db, sender) })
	auth.GET("/verify-email/:uid", func(c echo.Context) error { return verifyEmailHandler(c, db) })
	auth.Use(echojwt.WithConfig(config))
	auth.PATCH("/change-password", func(c echo.Context) error { return changePasswordHandler(c, db) })
	auth.POST("/forgot-password", func(c echo.Context) error { return forgotPasswordHandler(c, db) })
	auth.PATCH("/reset-password/:ptoken", func(c echo.Context) error { return resetPasswordHandler(c, db) })
	auth.POST("/logout", func(c echo.Context) error { return logoutHandler(c, db) })

}
