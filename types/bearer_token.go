package types

import "github.com/golang-jwt/jwt/v5"

// BeaerTokenClaims Bearer Token 声明
type BearerTokenClaims struct {
	jwt.RegisteredClaims
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
}
