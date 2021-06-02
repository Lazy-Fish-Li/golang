package user

import "golang.org/x/crypto/bcrypt"

type User struct {
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
}

func LoadTestUser() *User {
	//使用加密的“test”密码创建一个用户。

	//在真实的应用程序中，你可以通过特定的参数(电子邮件，用户名等)从数据库中加载用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("test1"), 8)
	return &User{Password: string(hashedPassword), Name: "Test user"}
}
