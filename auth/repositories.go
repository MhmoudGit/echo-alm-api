package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func getUserByEmail(db *gorm.DB, email string) (User, error) {
	var user User
	result := db.Where(User{Email: email}).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

// verify User password
func authinticateUser(db *gorm.DB, email, password string) (*User, error) {
	user, err := getUserByEmail(db, email)
	if err != nil {
		return nil, err
	}
	err = user.VerifyPassword(password)
	if err != nil {
		return nil, err
	}
	// Passwords match
	return &user, nil
}

func generateToken(userId, role string, duration time.Duration) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		userId,
		role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return t, nil
}
