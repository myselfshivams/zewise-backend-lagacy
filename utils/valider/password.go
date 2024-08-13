/*
Package valider - NekoBlog backend server data validers
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package valider

import "regexp"

// IsValidPassword 检查密码是否合法。
//
// 参数：
//   - password：密码
//
// 返回值：
//   - bool：如果密码合法，则返回true，否则返回false。
func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 32 {
		return false
	}

	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+={}\[\]:;'"<>,.?\/|\\~-]+$`)
	return re.MatchString(password)
}
