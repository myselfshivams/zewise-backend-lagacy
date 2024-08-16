/*
Package validers - NekoBlog backend server data validation.
This file is for username validation.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import "regexp"

// IsValidUsername 检查用户名是否合法。
//
// 参数：
//   - username：用户名
//
// 返回值：
//   - bool：如果用户名合法，则返回true，否则返回false。
func IsValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-z0-9_]+$`)
	return re.MatchString(username)
}
