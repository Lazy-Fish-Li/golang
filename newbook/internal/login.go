package internal

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

func UserLogin(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// 查询数据库
	var user User
	db, err := connectDB()
	if err != nil {
		return err
	}
	db.AutoMigrate(&storedUser{})
	db.First(&user, "username =?", u.Username)
	if user.Password == "" || user.Password != u.Password {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "UnAuthorized"})
	}
	//方法token
	claims := &jwtClaims{
		user.Username,
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
