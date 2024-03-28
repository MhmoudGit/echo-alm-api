package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	SUPERADMIN = "SuperAdmin"
	ADMIN      = "Admin"
	USER       = "User"
	// Regular expression for matching email addresses
    EMAILREGEX = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type JwtCustomClaims struct {
	UserID   string `json:"userId"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
	jwt.RegisteredClaims
}

// email verifications sender email and password
type EmailSender struct {
	Email    string
	Password string
}

type Login struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// user model in database
type User struct {
	ID        string    `gorm:"primary_key" json:"id"`
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

type UserCreate struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}

// serializes user creations from UserCreate form to insert to database
func (u *UserCreate) Serialize() *User {
	return &User{
		ID:       uuid.NewString(),
		Email:    u.Email,
		Password: u.Password,
		Role:     USER,
	}
}
