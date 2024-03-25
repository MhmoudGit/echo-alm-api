package auth

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	RestaurantID string `json:"restaurantId"`
	UserID       string `json:"userId"`
	Role         string `json:"role"`
	jwt.RegisteredClaims
}

type Login struct {
	Email    string `json:"email" form:"email" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
}
