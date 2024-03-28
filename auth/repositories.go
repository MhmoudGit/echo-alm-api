package auth

import (
	"fmt"
	"net/smtp"
	"regexp"
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
func authinticateUser(db *gorm.DB, email, password string) (User, error) {
	user, err := getUserByEmail(db, email)
	if err != nil {
		return user, err
	}
	err = user.VerifyPassword(password)
	if err != nil {
		return user, err
	}
	// Passwords match
	return user, nil
}

func generateToken(userId, role string, isActive bool, duration time.Duration, secret string) (string, error) {
	// Set custom claims
	claims := &JwtCustomClaims{
		userId,
		role,
		isActive,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func parseToken(tokenString string, secret string) (*JwtCustomClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the key used to sign the token
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func SendEmail(sender EmailSender, to []string, msg []byte) error {
	auth := smtp.PlainAuth("", sender.Email, sender.Password, "smtp.gmail.com")
	// Here we do it all: connect to our server, set up a message and send it
	err := smtp.SendMail("smtp.gmail.com:587", auth, sender.Email, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func CreateUser(db *gorm.DB, user *User) error {
	// Compile the regex pattern
	regexPattern := regexp.MustCompile(EMAILREGEX)
	ok := regexPattern.MatchString(user.Email)
	if !ok {
		return fmt.Errorf("invalid email")
	}
	err := user.HashPassword(user.Password)
	if err != nil {
		return err
	}
	// Create the User in the database
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
