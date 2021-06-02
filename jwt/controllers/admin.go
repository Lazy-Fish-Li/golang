package controllers

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Admin() echo.HandlerFunc {
	return func(c echo.Context) error {
        // 得到用户的cookie。
		userCookie, _ := c.Cookie("user")
		return c.String(http.StatusOK, fmt.Sprintf("Hi, %s! You have been authenticated!", userCookie.Value))
	}
}