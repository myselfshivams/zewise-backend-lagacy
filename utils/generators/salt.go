/*
Package generators - NekoBlog backend server generator utils
This file is for salt generator.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package generators

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateSalt 生成指定长度的盐。
//
// 参数：
//   - length：盐的长度
//
// 返回值：
//   - string：生成的盐。
//   - error：如果在生成盐的过程中发生错误，则返回相应的错误信息，否则返回nil。
func GenerateSalt(length int) (string, error) {
	// base64编码会将3个字节编码为4个字符
	numBytes := length / 4 * 3

	// 生成随机字节序列
	randomBytes := make([]byte, numBytes)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 将随机字节序列进行base64编码
	salt := base64.URLEncoding.EncodeToString(randomBytes)

	// 截取所需长度的字符串作为盐
	return salt[:length], nil
}
