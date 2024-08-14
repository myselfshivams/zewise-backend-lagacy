package controllers

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/models"
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

// NewProfileHandler 返回获取用户资料的处理函数。
//
// 返回值：
//   - fiber.Handler：新的获取用户资料的处理函数。
func (controller *UserController) NewProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 判断传入的查询参数是UID还是Username
		uidStr := ctx.Query("uid")
		username := ctx.Query("username")
		if uidStr == "" && username == "" {
			return ctx.Status(200).JSON(
				serializers.NewResponse(1, "parameter uid or username is required"),
			)
		}

		// 获取用户信息
		var (
			user *models.UserInfo
			err  error
		)
		if uidStr != "" {
			var uid uint64
			uid, err = strconv.ParseUint(uidStr, 10, 64)
			if err != nil {
				return ctx.Status(200).JSON(
					serializers.NewResponse(1, err.Error()),
				)
			}
			user, err = controller.userService.GetUserInfoByUID(uid)
		} else {
			user, err = controller.userService.GetUserInfoByUsername(username)
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(200).JSON(
					serializers.NewResponse(2, "user does not exist"),
				)
			}
			return ctx.Status(200).JSON(
				serializers.NewResponse(2, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(0, "", serializers.NewUserProfileData(user)),
		)
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
