package controllers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/useragent"
	"gorm.io/gorm"

	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/types"
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
				serializers.NewResponse(consts.PARAMETER_ERROR, "parameter uid or username is required"),
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
					serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
				)
			}
			user, err = controller.userService.GetUserInfoByUID(uid)
		} else {
			user, err = controller.userService.GetUserInfoByUsername(username)
		}
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(200).JSON(
					serializers.NewResponse(consts.PARAMETER_ERROR, "user does not exist"),
				)
			}
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "", serializers.NewUserProfileData(user)),
		)
	}
}

// NewRegisterHandler 返回注册用户的处理函数。
//
// 返回值：
//   - fiber.Handler：新的注册用户的处理函数。
func (controller *UserController) NewRegisterHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 注册用户
		err = controller.userService.RegisterUser(reqBody.Username, reqBody.Password)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

// NewLoginHandler 返回登陆用户的处理函数。
//
// 返回值：
//   - fiber.Handler：新的登陆用户的处理函数。
func (controller *UserController) NewLoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		reqBody := new(types.UserAuthBody)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 解析 UA
		userAgentString := ctx.Get("User-Agent")
		ua := useragent.New(userAgentString)
		browser, version := ua.Browser()
		var sb strings.Builder
		sb.WriteString(browser)
		sb.WriteString(" ")
		sb.WriteString(version)
		browserInfo := sb.String()
		os := ua.OSInfo().FullName

		// 登陆
		token, err := controller.userService.LoginUser(reqBody.Username, reqBody.Password, ctx.IP(), browserInfo, os)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		// 返回结果
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed", serializers.NewUserToken(token)),
		)
	}
}

// NewUploadAvatarHandler 返回上传头像的处理函数。
//
// 返回值：
//   - fiber.Handler：新的上传头像的处理函数。
func (controller *UserController) NewUploadAvatarHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// 获取Token Claims
		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		// 获取文件
		// 检查表单中文件的数量
		form, err := ctx.MultipartForm()
		if err != nil {
			return err
		}
		files := form.File["avatar"]
		if len(files) != 1 {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, "required 1 file, but got more or less"),
			)
		}
		fileHeader := files[0]

		// 保存头像
		err = controller.userService.UserUploadAvatar(claims.UID, fileHeader)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, err.Error()),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "succeed"),
		)
	}
}

//	NewUserUpdatePassword 修改密码的函数
//
// 返回值：
//   - fiber.Handler：新的上传头像的处理函数。
func (controller *UserController) NewUserUpdatePasswordHandler() fiber.Handler {
	type UserUpdatePassword struct {
		types.UserAuthBody
		NewPassword string `json:"new_password"`
	}

	return func(ctx *fiber.Ctx) error {
		// 解析请求体
		userUpdatePassword := new(UserUpdatePassword)
		err := ctx.BodyParser(userUpdatePassword)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		err = controller.userService.UserUpdatePassword(
			userUpdatePassword.Username,
			userUpdatePassword.Password,
			userUpdatePassword.NewPassword,
		)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.AUTH_ERROR, err.Error()),
			)
		}

		// 返回成功的 JSON 响应
		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "password updated successfully"),
		)
	}
}

func (controller *UserController) NewUserUpdateProfileHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		reqBody := new(types.UserUpdateProfile)
		err := ctx.BodyParser(reqBody)
		if err != nil {
			return ctx.Status(200).JSON(
				serializers.NewResponse(consts.PARAMETER_ERROR, err.Error()),
			)
		}

		claims := ctx.Locals("claims").(*types.BearerTokenClaims)

		err = controller.userService.UpdateUserInfo(claims.UID, reqBody)
		if err != nil {
			return ctx.Status(500).JSON(
				serializers.NewResponse(consts.SERVER_ERROR, "Failed to update profile"),
			)
		}

		return ctx.Status(200).JSON(
			serializers.NewResponse(consts.SUCCESS, "Profile updated successfully"),
		)
	}
}
