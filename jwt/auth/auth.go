package auth

import (
	"net/http"
	"time"

	"github.com/alexsergivan/blog-examples/authentication/user"
	"github.com/labstack/echo/v4"

	"github.com/dgrijalva/jwt-go"
)

const (
	accessTokenCookieName  = "access-token"
	refreshTokenCookieName = "refresh-token"
	//只是为了演示的目的，我在这里宣布了秘密。在实际应用程序中，您可能需要

	//从env变量中获取。
	jwtSecretKey        = "some-secret-key"
	jwtRefreshSecretKey = "some-refresh-secret-key"
)

func GetJWTSecret() string {
	return jwtSecretKey
}

func GetRefreshJWTSecret() string {
	return jwtRefreshSecretKey
}

//创建一个将被编码到JWT的结构体。

//我们添加jwt。StandardClaims作为嵌入式类型，提供过期时间等字段
type Claims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func GenerateTokensAndSetCookies(user *user.User, c echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}

	setTokenCookie(accessTokenCookieName, accessToken, exp, c)
	setUserCookie(user, exp, c)
	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(refreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func generateAccessToken(user *user.User) (string, time.Time, error) {
	// 声明令牌的过期时间
	expirationTime := time.Now().Add(1 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetJWTSecret()))
}

func generateRefreshToken(user *user.User) (string, time.Time, error) {
	// 声明令牌的过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(GetRefreshJWTSecret()))
}

func generateToken(user *user.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// 创建JWT声明，其中包括用户名和到期时间
	claims := &Claims{
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			// 在JWT中，过期时间表示为unix毫秒
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 使用用于签名的算法声明令牌和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 创建JWT字符串
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true

	c.SetCookie(cookie)
}

func setUserCookie(user *user.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Name
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

// 当用户试图访问受保护的路径时，将执行JWTErrorChecker。
func JWTErrorChecker(err error, c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}

// TokenRefresherMiddleware，它在访问令牌即将过期时刷新JWT令牌。
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 如果用户没有经过身份验证(在上下文中没有用户令牌数据)，那么什么也不要做。
		if c.Get("user") == nil {
			return next(c)
		}
		// 从上下文获取用户令牌。
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		//我们确保新令牌直到足够的时间过去才被发出

		//在这种情况下，新令牌只会在旧令牌在的情况下被发出

		//有效期15分钟。
		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 15*time.Minute {
			// 从cookie中获取刷新标记。
			rc, err := c.Cookie(refreshTokenCookieName)
			if err == nil && rc != nil {
				// 解析标记并检查它是否有效。
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(GetRefreshJWTSecret()), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// 如果一切正常，更新令牌。
					_ = GenerateTokensAndSetCookies(&user.User{
						Name: claims.Name,
					}, c)
				}
			}
		}

		return next(c)
	}
}
