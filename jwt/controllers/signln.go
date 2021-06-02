package controllers

import (
	"html/template"
	"net/http"
	"path"

	"github.com/alexsergivan/blog-examples/authentication/auth"
	"github.com/alexsergivan/blog-examples/authentication/user"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// SignInForm负责signIn表单的渲染。
func SignInForm() echo.HandlerFunc {
	return func(c echo.Context) error {
		fp := path.Join("templates", "signIn.html")
		tmpl, err := template.ParseFiles(fp)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		if err := tmpl.Execute(c.Response().Writer, nil); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return nil
	}
}

// SignIn将在SignInForm提交后执行。
func SignIn() echo.HandlerFunc {
	return func(c echo.Context) error {
		// Load  "test" user.
		storedUser := user.LoadTestUser()
		// 初始化一个新的User结构体。
		u := new(user.User)
		// 解析提交的数据并使用来自SignIn表单的数据填充User结构。
		if err := c.Bind(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// 将存储的哈希密码与收到的哈希版本的密码进行比较。
		if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(u.Password)); err != nil {
			// 如果两个密码不匹配，则返回401状态。
			return echo.NewHTTPError(http.StatusUnauthorized, "Password is incorrect")
		}
		// 如果密码正确，生成令牌并设置cookie。
		err := auth.GenerateTokensAndSetCookies(storedUser, c)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token is incorrect")
		}

		return c.Redirect(http.StatusMovedPermanently, "/admin")
	}
}
