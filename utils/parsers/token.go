/*
Package parsers - NekoBlog backend server data parsing utilities.
This file is for token parsing.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package parsers

import (
	"github.com/golang-jwt/jwt/v5"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// ParseToken 解析令牌。
//
// 参数：
//   - token：令牌字符串。
//
// 返回值：
//   - *BearerTokenClaims：令牌中的声明。
//   - error：如果在解析过程中发生错误，则返回相应的错误信息，否则返回nil。
func ParseToken(token string) (*types.BearerTokenClaims, error) {
	// 解析令牌
	claims := new(types.BearerTokenClaims)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SECRET), nil
	})

	// 返回结果
	return claims, err
}
