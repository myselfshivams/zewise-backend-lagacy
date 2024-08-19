/*
Package validers - NekoBlog backend server data validation.
This file is for token validation.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package validers

import (
	"time"

	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ValideTokenClaims 验证 Token 的 Claims。
//
// 参数
//   - claims：Token 的 Claims。
//
// 返回值
//   - bool：Token 的 Claims 是否有效。
func ValideTokenClaims(claims *types.BearerTokenClaims) bool {
	// 检测 Token 是否生效
	if time.Now().Unix() < claims.NotBefore.Unix() {
		return false
	}

	// 检测 Token 是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return false
	}

	return true
}
