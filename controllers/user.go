package controllers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/utils/serializers"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

// NewUserController 返回一个新的 UserController 实例。
//
// 返回值：
//   - *UserController：新的 UserController 实例。
func (factory *Factory) NewUserController() *UserController {
	return &UserController{
		userService: factory.serviceFactory.NewUserService(),
	}
}

// NewRegisterHandler 返回注册用户的处理函数。
//
// 返回值：
//   - fiber.Handler：新的注册用户的处理函数。
func (controller *UserController) NewRegisterHandler() fiber.Handler {
	// Body 请求结构
	type Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(Body)
		err := ctx.BodyParser(&reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(1, err.Error()),
			)
		}

		// 注册用户
		err = controller.userService.RegisterUser(reqBody.Username, reqBody.Password)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(2, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(0, "succeed"),
		)
	}
}
