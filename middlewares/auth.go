/*
Package middlewares - NekoBlog backend server middlewares.
This file is for token authentication middleware.
Copyright (c) [2024], Author(s):
- WhitePaper233<baizhiwp@gmail.com>
*/
package middlewares

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/parsers"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/validers"
)

// TokenAuthMiddleware 认证中间件
type TokenAuthMiddleware struct {
	userStore *stores.UserStore
}

// NewTokenAuthMiddleware 返回一个新的 AuthMiddleware 实例。
//
// 返回值
//   - *AuthMiddleware：新的 AuthMiddleware 实例。
func (factory *Factory) NewTokenAuthMiddleware() *TokenAuthMiddleware {
	return &TokenAuthMiddleware{userStore: factory.store.NewUserStore()}
}

// NewMiddleware Token 认证中间件
//
// 参数
//   - ctx：Fiber 上下文。
//
// 返回值
//   - error：错误
func (middleware *TokenAuthMiddleware) NewMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 从请求头中获取 Token
		token := ctx.Get("Authorization")
		if token == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "bearer token is required"),
			)
		}
		if len(token) < 7 || token[:7] != "Bearer " {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "bearer token is invalid"),
			)
		}
		token = token[7:]

		// 验证 Token
		claims, err := parsers.ParseToken(token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		// 检验 Token 是否在有效期内
		if !validers.ValideTokenClaims(claims) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is expired"),
			)
		}

		// 检验 Token 是否可用
		isAvaliable, err := middleware.userStore.IsUserTokenAvaliable(token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}
		if !isAvaliable {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is not avaliable"),
			)
		}

		// 将 claims 信息存入 ctx.Locals 中
		ctx.Locals("claims", claims)

		return ctx.Next()
	}
}
