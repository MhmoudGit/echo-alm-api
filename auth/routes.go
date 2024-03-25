package auth

import (
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func AuthRoutes(e *echo.Echo, config echojwt.Config, Db *gorm.DB) {
	auth := e.Group("/auth")

	auth.POST("/token", func(c echo.Context) error { return tokenHandler(c, Db) })
	auth.POST("/refresh-token", func(c echo.Context) error { return refreshTokenHandler(c, Db) })
	auth.POST("/signup", func(c echo.Context) error { return signupHandler(c, Db) })
	auth.POST("/verify-email/:uid", func(c echo.Context) error { return verifyEmailHandler(c, Db) })
	auth.POST("/change-password", func(c echo.Context) error { return changePasswordHandler(c, Db) })
	auth.POST("/forgot-password", func(c echo.Context) error { return forgotPasswordHandler(c, Db) })
	auth.POST("/reset-password/:ptoken", func(c echo.Context) error { return resetPasswordHandler(c, Db) })
	auth.POST("/logout", func(c echo.Context) error { return logoutHandler(c, Db) })

}
