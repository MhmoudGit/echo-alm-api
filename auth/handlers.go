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
func tokenHandler(c echo.Context, db *gorm.DB, secret string) error {
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

	// check if user is active then generate the access token and the refresh token
	var accessToken string
	var refreshToken string
	if user.IsActive {
		accessToken, err = generateToken(user.ID, user.Role, user.IsActive, time.Minute*10, secret)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "unable to generate token, check user activity")
		}
		refreshToken, err = generateToken(user.ID, user.Role, user.IsActive, time.Hour*24, secret)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "unable to generate refresh token")
		}
	}

	// Create a new cookie
	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.SameSite = http.SameSiteLaxMode
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, echo.Map{
		"token": accessToken,
	})
}

// @Tags		auth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/refresh-token [post]
func refreshTokenHandler(c echo.Context, secret string) error {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return echo.ErrUnauthorized
	}
	claims, err := parseToken(cookie.Value, secret)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	var accessToken string
	if claims.IsActive {
		accessToken, err = generateToken(claims.UserID, claims.Role, claims.IsActive, time.Minute*10, secret)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "unable to generate token, check user activity")
		}
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": accessToken,
	})
}

// @Tags		auth
// @Produce	json
// @Param		user	formData	auth.UserCreate	true	"new user data"
// @Success	200		{string}	{string}		"Successful Response"
// @Router		/auth/signup [post]
func signupHandler(c echo.Context, db *gorm.DB, sender EmailSender) error {
	var userData UserCreate
	err := c.Bind(&userData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err = c.Validate(&userData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	newUser := userData.Serialize()

	// check if email exists in db
	_, err = getUserByEmail(db, newUser.Email)
	if err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "change your email")
	}

	// create user in db
	err = CreateUser(db, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	SendEmail(sender, []string{newUser.Email}, []byte("Hello world"))

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user registered successfully",
	})
}

// @Tags		auth
// @security	ApiKeyAuth
// @Produce	json
// @Success	200	{string}	string	"Successful Response"
// @Router		/auth/verify-email/{uid} [get]
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
