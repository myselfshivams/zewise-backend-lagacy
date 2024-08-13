/*
Package encryptor - NekoBlog backend server data encryptors
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package encryptor

import "golang.org/x/crypto/bcrypt"

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
	passwordWithSalt := append([]byte(password), []byte(salt)...)

	hashedPassword, err := bcrypt.GenerateFromPassword(
		passwordWithSalt,
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}
