package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JwtCustomClaims struct {
	UserID   string `json:"userId"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
	jwt.RegisteredClaims
}

type Login struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

type User struct {
	ID        string    `gorm:"not null;index;unique" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `gorm:"not null;index;unique" json:"email" form:"email"`
	Password  string    `gorm:"not null" json:"-" form:"password"`
	Role      string    `gorm:"not null" json:"role" form:"role"`
	IsActive  bool      `gorm:"not null;default:true" json:"isActive" form:"isActive"`
}

// Verify Password.
func (u *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err
}

// Verify Activity.
func (u *User) VerifyActivity(isActive bool) bool {
	return u.IsActive
}

// HashPassword securely hashes the provided password and sets it in the Password field.
func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
