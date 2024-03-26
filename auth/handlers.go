package auth

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// @Tags		auth
// @Accept		x-www-form-urlencoded
// @Produce	json
// @Param		login	formData	auth.Login	true	"login data"
// @Success	200		{string}	string		"Successful Response"
// @Router		/auth/token [post]
func tokenHandler(c echo.Context, db *gorm.DB) error {
	var logInfo Login
	err := c.Bind(&logInfo)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(&logInfo); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	user, err := authinticateUser(db, logInfo.Email, logInfo.Password)
	if err != nil {
		return echo.ErrUnauthorized
	}

	token, err := generateToken(user.ID, user.Role, time.Minute*10)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Unable to get token")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token, "user": &user,
	})
}

// @Tags		auth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/refresh-token [post]
func refreshTokenHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @Produce	json
// @Param		login	formData	auth.Login	true	"new user data"
// @Success	200		{string}	auth.Login	"Successful Response"
// @Router		/auth/signup [post]
func signupHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/verify-email/{uid} [patch]
func verifyEmailHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/change-password [patch]
func changePasswordHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Param		email	formData	string	true	"registered email"
// @Success	200		{string}	string	"Successful Response"
// @Router		/auth/forgot-password [post]
func forgotPasswordHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Param		login	formData	auth.Login	true	"login data"
// @Success	200		{string}	string		"Successful Response"
// @Router		/auth/reset-password/{ptoken} [patch]
func resetPasswordHandler(c echo.Context, db *gorm.DB) error {
	return nil
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/logout [post]
func logoutHandler(c echo.Context, db *gorm.DB) error {
	return nil
}
