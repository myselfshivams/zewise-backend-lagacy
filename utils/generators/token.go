/*
Package generator - NekoBlog backend server data generators
This file is for token generator.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package generators

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
)

// GenerateToken 生成一个新的令牌。
//
// 参数：
//   - uid：用户ID
//   - username：用户名
//
// 返回值：
//   - string：新的令牌。
//   - error：如果在生成过程中发生错误，则返回相应的错误信息，否则返回nil。
func GenerateToken(uid uint64, username string) (string, *types.BearerTokenClaims, error) {
	// 构造 Token 的 Claims
	claims := &types.BearerTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(consts.TOKEN_EXPIRE_DURATION * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    consts.TOKEN_ISSUER,
			Subject:   "BearerToken",
			ID:        uuid.New().String(),
		},
		UID:      uid,
		Username: username,
	}

	// 生成 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 Token
	tokenString, err := token.SignedString([]byte(consts.TOKEN_SECRET))

	// 返回 Token 和 Claims
	return tokenString, claims, err
}
