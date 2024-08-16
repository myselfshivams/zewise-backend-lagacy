/*
Package type - NekoBlog backend server types.
This file is for token related types.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package types

import "github.com/golang-jwt/jwt/v5"

// BeaerTokenClaims Bearer Token 声明
type BearerTokenClaims struct {
	jwt.RegisteredClaims
	UID      uint64 `json:"uid"`
	Username string `json:"username"`
}
