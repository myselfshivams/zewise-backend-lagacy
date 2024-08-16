package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/parser"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/valider"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	userStore *stores.UserStore
}

// NewAuthMiddleware 返回一个新的 AuthMiddleware 实例。
//
// 返回值
//   - *AuthMiddleware：新的 AuthMiddleware 实例。
func (factory *Factory) NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{userStore: factory.store.NewUserStore()}
}

// NewTokenAuth Token 认证中间件
//
// 参数
//   - ctx：Fiber 上下文。
//
// 返回值
//   - error：错误
func (middleware *AuthMiddleware) NewTokenAuth() fiber.Handler {
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
		claims, err := parser.ParseToken(token)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		// 检验 Token 是否在有效期内
		if !valider.ValideTokenClaims(claims) {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, "bearer token is expired"),
			)
		}

		// 检验 Token 是否可用
		isAvaliable, err := middleware.userStore.IsTokenAvaliable(token)
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
