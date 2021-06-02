package main

import (
	"net/http"
	"newbook/internal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/storyicon/grbac"
)

func QueryRolesByHeaders(header http.Header) (roles []string, err error) {
	// 在这里实现你的逻辑
	// ...
	// 这个逻辑可能是从请求的Headers中获取token，并且根据token从数据库中查询用户的相应角色。

	return roles, err
}

func Authentication() echo.MiddlewareFunc {
	rbac, err := grbac.New(grbac.WithJSON("config.jsom", time.Minute*10))
	if err != nil {
		panic(err)
	}
	return func(echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			roles, err := QueryRolesByHeaders(c.Request().Header)
			if err != nil {
				c.NoContent(http.StatusInternalServerError)
				return nil
			}
			state, err := rbac.IsRequestGranted(c.Request(), roles)
			if err != nil {
				c.NoContent(http.StatusInternalServerError)
				return nil
			}
			if state.IsGranted() {
				return nil
			}
			c.NoContent(http.StatusUnauthorized)
			return nil
		}
	}
}

func main() {
	e := echo.New()
	e.POST("/login", internal.UserLogin)
	e.Logger.Fatal(e.Start(":1323"))
}
