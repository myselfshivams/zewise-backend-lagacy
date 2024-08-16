/*
Package encryptors - NekoBlog backend server data encryptors.
This file is for password encryptors.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package encryptors

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 生成哈希密码。
//
// 参数：
//   - password：密码
//   - salt：盐
//
// 返回值：
//   - string：生成的哈希密码。
//   - error：如果在生成哈希密码的过程中发生错误，则返回相应的错误信息，否则返回nil。
func HashPassword(password string, salt string) (string, error) {
	// 将密码和盐拼接在一起
	passwordWithSalt := append([]byte(password), []byte(salt)...)

	// 生成哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword(
		passwordWithSalt,
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	// 返回哈希密码
	return string(hashedPassword), nil
}

// CompareHashPassword 比较哈希密码。
//
// 参数：
//   - hashedPassword：哈希密码
//   - password：密码
//   - salt：盐
//
// 返回值：
//   - error：如果密码匹配，则返回nil，否则返回相应的错误信息。
func CompareHashPassword(hashedPassword string, password string, salt string) error {
	// 将密码和盐拼接在一起
	passwordWithSalt := append([]byte(password), []byte(salt)...)

	// 比较哈希密码
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), passwordWithSalt)
}
