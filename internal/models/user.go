package models

import (
	"database/sql"

	"github.com/golang-jwt/jwt"
)

type User struct {
	ID                      uint         `json:"id" gorm:"primaryKey"`
	Username                string       `json:"username"`
	Password                string       `json:"-"`
	UnlimitedSwipeExpiredAt sql.NullTime `json:"-"`
}

type AuthRequest struct {
	Username string `json:"username" validate:"required,min=4,max=20"`
	Password string `json:"password" validate:"required,min=6"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}
