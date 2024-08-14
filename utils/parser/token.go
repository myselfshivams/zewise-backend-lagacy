package parser

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
	claims, err := jwt.ParseWithClaims(token, &types.BearerTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(consts.TOKEN_SECRET), nil
	})
	return claims.Claims.(*types.BearerTokenClaims), err
}
