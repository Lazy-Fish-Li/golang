package internal

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username" gorm:"username"`
	Password string `json:"password" gorm:"password"`
	Role     string `json:"role" gorm:"role"`
}

type jwtClaims struct {
	Username string `json:"username" gorm:"username"`
	Role     string `json:"role" gorm:"role"`
	jwt.StandardClaims
}

type storedUser struct {
	gorm.Model
	Role string `json:"role" gorm:"role"`
}
